package tx

import (
	"time"

	"github.com/deadsy/libusb"
)

type TxEngine struct {
	ctx    libusb.Context
	handle libusb.Device_Handle

	opened      bool
	lastCommand time.Time
}

func NewTxEngine() (*TxEngine, error) {
	engine := TxEngine{}

	if err := libusb.Init(&engine.ctx); err != nil {
		return nil, err
	}

	return &engine, nil
}

func (this *TxEngine) Open() error {
	if this.opened {
		panic("Already opened")
	}

	this.handle = libusb.Open_Device_With_VID_PID(this.ctx, tx_VendorID, tx_ProductID)
	var driver_active bool
	var err error
	if driver_active, err = libusb.Kernel_Driver_Active(this.handle, tx_Interface); err != nil {
		libusb.Close(this.handle)
		return err
	} else if driver_active {
		libusb.Detach_Kernel_Driver(this.handle, tx_Interface)
	}
	if err := libusb.Set_Configuration(this.handle, tx_Config); err != nil {
		libusb.Close(this.handle)
		return err
	}
	if err := libusb.Claim_Interface(this.handle, tx_Interface); err != nil {
		libusb.Close(this.handle)
		return err
	}
	this.opened = true

	return nil
}

func (this *TxEngine) Close() {
	if this.opened {
		libusb.Attach_Kernel_Driver(this.handle, tx_Interface)
		libusb.Close(this.handle)
	}
}

func (this *TxEngine) Exit() {
	libusb.Exit(this.ctx)
}

func (this *TxEngine) Write(c Command) error {
	sleepDelay := this.lastCommand.Add(tx_Delay).Sub(time.Now())
	if sleepDelay > 0 {
		time.Sleep(sleepDelay)
	}
	_, err := libusb.Control_Transfer(this.handle, tx_ReqType|libusb.ENDPOINT_OUT, 0x9, 0x300, 0, c.Data(nil), 200)
	this.lastCommand = time.Now()
	return err
}
