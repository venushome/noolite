package tx

const (
	TX_CMD_TURN_OFF     byte = 0
	TX_CMD_DIM_DOWN     byte = 1
	TX_CMD_TURN_ON      byte = 2
	TX_CMD_DIM_UP       byte = 3
	TX_CMD_TURN_CHANGE  byte = 4
	TX_CMD_DIM_CHANGE   byte = 5
	TX_CMD_DIM_SET      byte = 6
	TX_CMD_SCENE_RUN    byte = 7
	TX_CMD_SCENE_SAVE   byte = 8
	TX_CMD_CLEAR_ADDR   byte = 9
	TX_CMD_DIM_STOP     byte = 10
	TX_CMD_WANT_BIND    byte = 15
	TX_CMD_RAINBOW      byte = 16
	TX_CMD_CHANGE_COLOR byte = 17
	TX_CMD_WORK_MODE    byte = 18
	TX_CMD_WORK_SPEED   byte = 19
)

var (
	zeroBytes        []byte          = []byte{0x0, 0x0, 0x0, 0x0}
	commandFormatMap map[byte][]byte = map[byte]([]byte){
		TX_CMD_TURN_OFF:     zeroBytes,
		TX_CMD_DIM_DOWN:     zeroBytes,
		TX_CMD_TURN_ON:      zeroBytes,
		TX_CMD_DIM_UP:       zeroBytes,
		TX_CMD_TURN_CHANGE:  zeroBytes,
		TX_CMD_DIM_CHANGE:   zeroBytes,
		TX_CMD_DIM_SET:      []byte{0x1, 0x3, 0x3, 0x3},
		TX_CMD_SCENE_RUN:    zeroBytes,
		TX_CMD_SCENE_SAVE:   zeroBytes,
		TX_CMD_CLEAR_ADDR:   zeroBytes,
		TX_CMD_DIM_STOP:     zeroBytes,
		TX_CMD_WANT_BIND:    zeroBytes,
		TX_CMD_RAINBOW:      []byte{0x4, 0x4, 0x4, 0x4},
		TX_CMD_CHANGE_COLOR: []byte{0x4, 0x4, 0x4, 0x4},
		TX_CMD_WORK_MODE:    []byte{0x4, 0x4, 0x4, 0x4},
		TX_CMD_WORK_SPEED:   []byte{0x4, 0x4, 0x4, 0x4},
	}
)

type Command struct {
	Type    byte
	Channel byte
	value   []byte
}

func (this Command) Data(buf []byte) (data []byte) {
	if buf == nil {
		data = make([]byte, 0, 8)
	} else {
		data = buf[:0]
	}

	data = append(data,
		0x30, // 0x30 - 1 repeat + 1000b/sec
		this.Type,
		commandFormatMap[this.Type][len(this.value)],
		0x0,
		this.Channel,
	)

	data = append(data, this.value...)

	for cnt := len(data); cnt < 8; cnt++ {
		data = append(data, 0x0)
	}
	return data
}

func (this *Command) SetDimLevel(level byte) {
	this.value = []byte{level}
}

func (this *Command) SetRGB(r, g, b byte) {
	this.value = []byte{r, g, b}
}
