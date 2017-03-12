package tx

import (
	"github.com/deadsy/libusb"
	"time"
)

const (
	tx_VendorID  = 5824
	tx_ProductID = 1503
	tx_Interface = 0
	tx_Config    = 1
	tx_ReqType   = libusb.REQUEST_TYPE_CLASS | libusb.RECIPIENT_INTERFACE
	tx_Delay     = 150 * time.Millisecond
)
