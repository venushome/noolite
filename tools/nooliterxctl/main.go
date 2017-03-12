package main

import (
	"fmt"

	"log"
	"strconv"
	"strings"
	"time"

	"github.com/chzyer/readline"
	"github.com/venushome/noolite/rx"
)

var (
	engine       *rx.RxEngine
	commandNames map[byte]string = map[byte]string{
		rx.RX_RESP_TURN_OFF:    "turn off",
		rx.RX_RESP_DIM_DOWN:    "dim light down",
		rx.RX_RESP_TURN_ON:     "turn on",
		rx.RX_RESP_DIM_UP:      "dim light up",
		rx.RX_RESP_TURN_CHANGE: "change state",
		rx.RX_RESP_DIM_CHANGE:  "change dim direction",
		rx.RX_RESP_DIM_SET:     "set dim level",
		rx.RX_RESP_SCENE_RUN:   "run saved scene",
		rx.RX_RESP_SCENE_SAVE:  "save scene",
		rx.RX_RESP_CLEAR_ADDR:  "clear address",
		rx.RX_RESP_DIM_STOP:    "stop light dim",
		rx.RX_RESP_WANT_BIND:   "want to bind",
		rx.RX_RESP_RAINBOW:     "continious color change",
		rx.RX_RESP_SET_COLOR:   "set color",
		rx.RX_RESP_WORK_MODE:   "change work mode",
		rx.RX_RESP_WORK_SPEED:  "change work speed",
		rx.RX_RESP_LOW_BATERY:  "low batery",
		rx.RX_RESP_SENSE_INFO:  "sensor information",
	}
)

func openDev() {
	if e, err := rx.NewRxEngine(); err != nil {
		log.Printf("Device open filed: %v \n", err)
	} else {
		engine = e
	}
	err := engine.Open()
	if err != nil {
		log.Printf("Device open filed: %v \n", err)
		return
	}
	log.Println("Device opened")
	go dumpReads()
}

func formatResponse(r rx.Response) string {
	if r.Empty() {
		return "no data"
	}
	channel := r.Channel()
	commandName := commandNames[r.Command()]
	dataLen := r.DataLen()
	data := r.Data()

	cmdString := fmt.Sprintf("cmd: %s for channel: %d", commandName, channel)
	if dataLen > 0 {
		cmdString += fmt.Sprintf(" with %d bytes of data: %X", dataLen, data)
	}

	return cmdString
}

func dumpReads() {
	for {
		if engine == nil {
			return
		}
		resp, err := engine.Read(time.Second)
		if err != nil {
			log.Printf("Read error: ", err)
			return
		}
		if len(resp) > 0 {
			log.Printf("Got %s\n", formatResponse(resp))
		}
	}
}
func filterInput(r rune) (rune, bool) {
	switch r {
	// block CtrlZ feature
	case readline.CharCtrlZ:
		return r, false
	}
	return r, true
}

func main() {
	l, err := readline.NewEx(&readline.Config{
		Prompt:              "\033[31m»\033[0m ",
		HistoryFile:         "/tmp/readline.tmp",
		InterruptPrompt:     "^C",
		EOFPrompt:           "exit",
		HistorySearchFold:   true,
		FuncFilterInputRune: filterInput,
	})
	if err != nil {
		panic(err)
	}
	defer l.Close()

	log.SetOutput(l.Stderr())
	for {
		next := true
		line, err := l.Readline()
		if err == readline.ErrInterrupt {
			if len(line) == 0 {
				break
			} else {
				continue
			}
		}

		line = strings.TrimSpace(line)
		switch {
		case line == "open":
			openDev()

		case strings.HasPrefix(line, "bind "):
			line := strings.TrimSpace(line[4:])
			chanel, err := strconv.Atoi(line)
			if err != nil || chanel < 0 || chanel > 63 {
				log.Println("Bad channel for bind")
			} else {
				engine.Write(rx.Command{Type: rx.RX_CMD_BIND, Channel: byte(chanel)})
			}
		case strings.HasPrefix(line, "cancle "):
			line := strings.TrimSpace(line[7:])
			chanel, err := strconv.Atoi(line)
			if err != nil || chanel < 0 || chanel > 63 {
				log.Println("Bad channel for bind")
			} else {
				engine.Write(rx.Command{Type: rx.RX_CMD_CANCLE_BIND, Channel: byte(chanel)})
			}
		case strings.HasPrefix(line, "unbind "):
			line := strings.TrimSpace(line[7:])
			chanel, err := strconv.Atoi(line)
			if err != nil || chanel < 0 || chanel > 63 {
				log.Println("Bad channel for bind")
			} else {
				engine.Write(rx.Command{Type: rx.RX_CMD_CLEAR, Channel: byte(chanel)})
			}
		case line == "reset":
			engine.Write(rx.Command{Type: rx.RX_CMD_CLEAR_ALL})
		case line == "close":
			engine.Close()
			engine.Exit()
			engine = nil
		case line == "exit":
			next = false
		default:
			log.Println("you said:", strconv.Quote(line))
		}

		if next {
			continue
		}
		if engine != nil {
			engine.Close()
			engine.Exit()
		}
		log.Println("Bye")
		break

	}
}
