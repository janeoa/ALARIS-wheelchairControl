package main

import (
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"reflect"
	"time"

	"github.com/fatih/color"
)

func rx(f io.ReadWriteCloser) {
	// buff := make([]byte, 100)
	for {
		for {
			buf := make([]byte, 32)
			n, err := f.Read(buf)
			if err != nil {
				if err != io.EOF {
					fmt.Println("Error reading from serial port: ", err)
					f.Close()
					Arduino.isConnected = false
					break
				}
			} else {
				buf = buf[:n]
				color.Yellow("%s", string(buf))
				color.Cyan("%s", hex.Dump(buf))
			}
		}
		for !Arduino.isConnected {
			time.Sleep(time.Second * 1)
			color.Yellow("retrying Arduino connection")
		}
	}
}

func sendSingleCommand(port io.ReadWriteCloser, command byte) {
	b := commandPack{command, 0, 0}
	EasyTransferSend(port, b)
}

func sendCommand(port io.ReadWriteCloser, action byte, x_axis float32, y_axis float32) {
	x_byte := byte(x_axis + 100)
	y_byte := byte(y_axis + 100)
	b := commandPack{action, x_byte, y_byte}
	EasyTransferSend(port, b)
}

type commandPack struct {
	Action  byte
	Ch_name byte `uri:"ch_name"  binding:"number"`
	Value   byte `uri:"value" binding:"number"`
}

func EasyTransferSend(port io.ReadWriteCloser, in commandPack) {

	size := reflect.TypeOf(in).Size()
	CS := byte(size)
	toOut := []byte{0x06, 0x85}
	toOut = append(toOut, byte(size))

	toOut = append(toOut, in.Action)
	CS ^= in.Action
	toOut = append(toOut, in.Ch_name)
	CS ^= in.Ch_name
	toOut = append(toOut, in.Value)
	CS ^= in.Value

	toOut = append(toOut, CS)
	// if printUARTlogs {
	color.Cyan("Writing %v, as %v bytes using EasyTransfer\n", in, toOut)
	// }
	if Arduino.isConnected {
		_, err := port.Write(toOut)
		if err != nil {
			log.Fatalf("port.Write: %v", err)
		}
	} else {
		color.Red("EasyTranfer: no device connected")
	}
}
