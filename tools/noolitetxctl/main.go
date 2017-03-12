package main

import (
	"errors"
	"log"
	"strconv"
	"strings"

	"github.com/chzyer/readline"
	"github.com/venushome/noolite/tx"
)

var (
	engine      *tx.TxEngine
	mapCommands map[string]byte = map[string]byte{
		"off":        tx.TX_CMD_TURN_OFF,
		"dim-down":   tx.TX_CMD_DIM_DOWN,
		"on":         tx.TX_CMD_TURN_ON,
		"dim-up":     tx.TX_CMD_DIM_UP,
		"turn":       tx.TX_CMD_TURN_CHANGE,
		"dim-change": tx.TX_CMD_DIM_CHANGE,
		"dim":        tx.TX_CMD_DIM_SET,
		"run":        tx.TX_CMD_SCENE_RUN,
		"save":       tx.TX_CMD_SCENE_SAVE,
		"unbind":     tx.TX_CMD_CLEAR_ADDR,
		"dim-stop":   tx.TX_CMD_DIM_STOP,
		"bind":       tx.TX_CMD_WANT_BIND,
		"rainbow":    tx.TX_CMD_RAINBOW,
		"color":      tx.TX_CMD_CHANGE_COLOR,
		"mode":       tx.TX_CMD_WORK_MODE,
		"speed":      tx.TX_CMD_WORK_SPEED,
	}
)

func openDev() {
	if e, err := tx.NewTxEngine(); err != nil {
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
}

func parseRGB(s string, l, r, g, b *byte, is_rgb *bool) error {
	if strings.ContainsRune(s, rune(':')) {
		*is_rgb = true
		parts := strings.Split(s, ":")
		if len(parts) != 3 {
			return errors.New("Bad rgb value format")
		}
		ir, _ := strconv.Atoi(parts[0])
		ig, _ := strconv.Atoi(parts[1])
		ib, _ := strconv.Atoi(parts[2])
		*r = byte(ir)
		*g = byte(ig)
		*b = byte(ib)
		return nil
	}
	*is_rgb = false
	il, _ := strconv.Atoi(s)
	*l = byte(il)
	return nil
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
		Prompt:          "\033[31m»\033[0m ",
		HistoryFile:     "/tmp/readline.tmp",
		InterruptPrompt: "^C",
		EOFPrompt:       "exit",

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
		fields := strings.Fields(line)
		cmd, cmd_ok := mapCommands[fields[0]]
		if len(fields) == 0 {
			continue
		}
		switch {
		case line == "open":
			openDev()

		case cmd_ok:
			if len(fields) < 2 {
				log.Println("Not enouth arguments")
			}
			chanel, err := strconv.Atoi(fields[1])
			if err != nil || chanel < 0 || chanel > 63 {
				log.Println("Bad channel for bind")
				break
			}
			command := tx.Command{Type: cmd, Channel: byte(chanel)}
			switch cmd {
			case tx.TX_CMD_DIM_SET:
				if len(fields) == 3 {
					var r, g, b, l byte
					var is_rgb bool
					err := parseRGB(fields[2], &l, &r, &g, &b, &is_rgb)
					if err != nil {
						log.Printf("Error: %v \n", err)
						continue
					}
					if is_rgb {
						command.SetRGB(r, g, b)
					} else {
						command.SetDimLevel(l)
					}
				} else {
					log.Println("Not enouth arguments")
				}
			case tx.TX_CMD_RAINBOW, tx.TX_CMD_CHANGE_COLOR, tx.TX_CMD_WORK_MODE, tx.TX_CMD_WORK_SPEED:
				if len(fields) == 3 {
					var r, g, b byte
					var is_rgb bool
					if err != nil {
						log.Printf("Error: %v \n", err)
					} else if !is_rgb {
						log.Println("Need R:G:B value")
					} else {
						command.SetRGB(r, g, b)
					}
				} else {
					log.Println("Not enouth arguments")
				}
			}
			engine.Write(command)
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
