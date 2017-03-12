package rx

import (
	"bytes"
)

const (
	RX_RESP_TURN_OFF    byte = 0
	RX_RESP_DIM_DOWN    byte = 1
	RX_RESP_TURN_ON     byte = 2
	RX_RESP_DIM_UP      byte = 3
	RX_RESP_TURN_CHANGE byte = 4
	RX_RESP_DIM_CHANGE  byte = 5
	RX_RESP_DIM_SET     byte = 6
	RX_RESP_SCENE_RUN   byte = 7
	RX_RESP_SCENE_SAVE  byte = 8
	RX_RESP_CLEAR_ADDR  byte = 9
	RX_RESP_DIM_STOP    byte = 10
	RX_RESP_WANT_BIND   byte = 15
	RX_RESP_RAINBOW     byte = 16
	RX_RESP_SET_COLOR   byte = 17
	RX_RESP_WORK_MODE   byte = 18
	RX_RESP_WORK_SPEED  byte = 19
	RX_RESP_LOW_BATERY  byte = 20
	RX_RESP_SENSE_INFO  byte = 21
)

type Response []byte

func (this Response) Empty() bool {
	return bytes.Equal(this, []byte{0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0})
}

// togl value for new response mark
func (this Response) Togl() byte {
	return this[0] & 63
}

// check if new response received
func (this Response) NewResponse(prevTogl byte) (is_new bool, togl byte) {
	togl = this.Togl()
	is_new = togl != prevTogl
	return
}

// response channel (0-63)
func (this Response) Channel() uint8 {
	return uint8(this[1] & 63)
}

// response command type
func (this Response) Command() byte {
	return this[2]
}

// length of data in response
func (this Response) DataLen() int {
	switch this[3] {
	case byte(1):
		return 1
	case byte(2):
		return 2
	case byte(3):
		return 4
	}
	return 0
}

// response data
func (this Response) Data() []byte {
	return this[4 : 4+this.DataLen()]
}
