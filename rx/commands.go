package rx

const (
	RX_CMD_BIND        byte = 1
	RX_CMD_CANCLE_BIND byte = 2
	RX_CMD_CLEAR       byte = 3
	RX_CMD_CLEAR_ALL   byte = 4
)

type Command struct {
	Type    byte
	Channel byte
}

func (this Command) Data(buf []byte) (data []byte) {
	if buf == nil {
		data = make([]byte, 0, 8)
	} else {
		data = buf[:0]
	}
	return append(data, this.Type, this.Channel, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0)
}
