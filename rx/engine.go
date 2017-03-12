package rx

import (
	"sync"
	"time"

	"github.com/deadsy/libusb"
)

type RxEngine struct {
	togl   byte
	ctx    libusb.Context
	handle libusb.Device_Handle

	reads  chan Response
	writes chan Command

	opened bool
	stop   bool
	stoped chan bool
	err    error
	lock   sync.RWMutex
}

func NewRxEngine() (*RxEngine, error) {
	engine := RxEngine{
		togl: 255,
	}

	if err := libusb.Init(&engine.ctx); err != nil {
		return nil, err
	}

	return &engine, nil
}

func (this *RxEngine) Open() error {
	if this.opened {
		panic("Already opened")
	}

	this.handle = libusb.Open_Device_With_VID_PID(this.ctx, rx_VendorID, rx_ProductID)
	var driver_active bool
	var err error
	if driver_active, err = libusb.Kernel_Driver_Active(this.handle, rx_Interface); err != nil {
		libusb.Close(this.handle)
		return err
	} else if driver_active {
		libusb.Detach_Kernel_Driver(this.handle, rx_Interface)
	}
	if err := libusb.Set_Configuration(this.handle, rx_Config); err != nil {
		libusb.Close(this.handle)
		return err
	}
	if err := libusb.Claim_Interface(this.handle, rx_Interface); err != nil {
		libusb.Close(this.handle)
		return err
	}
	this.reads = make(chan Response, 100)
	this.writes = make(chan Command)
	this.stop = false
	this.stoped = make(chan bool)
	this.opened = true
	this.err = nil

	go this.ioCycle()

	return nil
}

func (this *RxEngine) Close() {
	if this.opened {
		this.stopCycle()
		<-this.stoped

		close(this.reads)
		close(this.writes)

		libusb.Attach_Kernel_Driver(this.handle, rx_Interface)
		libusb.Close(this.handle)
	}
}

func (this *RxEngine) Error() error {
	this.lock.RLock()
	defer this.lock.RUnlock()
	return this.err
}

func (this *RxEngine) Exit() {
	libusb.Exit(this.ctx)
}

func (this *RxEngine) stopCycle() {
	this.lock.Lock()
	defer this.lock.Unlock()
	this.stop = true
}

func (this *RxEngine) setError(err error) {
	this.lock.Lock()
	defer this.lock.Unlock()
	this.err = err
}

func (this *RxEngine) ioCycle() {
	checkCycle := func() bool {
		this.lock.RLock()
		defer this.lock.RUnlock()
		return this.stop
	}
	for {
		if checkCycle() {
			this.stoped <- true
			return
		}
		c := Command{}
		select {
		case c = <-this.writes:
			libusb.Control_Transfer(this.handle, rx_ReqType|libusb.ENDPOINT_OUT, 0x9, 0x300, 0, c.Data(nil), 200)
			time.Sleep(100 * time.Millisecond)
			continue
		default:
		}
		buf := make([]byte, 8)
		data, err := libusb.Control_Transfer(this.handle, rx_ReqType|libusb.ENDPOINT_IN, 0x9, 0x300, 0, buf, 250)
		if err != nil {
			this.setError(err)
			this.stoped <- true
			return
		} else {
			resp := Response(data)
			var is_new bool
			is_new, this.togl = resp.NewResponse(this.togl)
			if is_new {
				this.reads <- resp
			}
		}
		time.Sleep(50 * time.Millisecond)
	}
}

func (this *RxEngine) Read(timeout time.Duration) (Response, error) {
	if err := this.Error(); err != nil {
		return nil, err
	}
	select {
	case resp := <-this.reads:
		return resp, nil
	case <-time.After(timeout):
		if err := this.Error(); err != nil {
			return nil, err
		}
	}
	return nil, nil
}

func (this *RxEngine) Write(c Command) error {
	if err := this.Error(); err != nil {
		return err
	}
	this.writes <- c
	return nil
}
