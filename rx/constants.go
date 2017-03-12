package rx

import (
	"github.com/deadsy/libusb"
)

const (
	rx_VendorID  = 5824
	rx_ProductID = 1500
	rx_Interface = 0
	rx_Config    = 1
	rx_ReqType   = libusb.REQUEST_TYPE_CLASS | libusb.RECIPIENT_INTERFACE
)
