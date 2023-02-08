package main

import (
	"bytes"
	"fmt"
	"log"

	"github.com/digitalbitbox/bitbox02-api-go/communication/u2fhid"
	"github.com/digitalbitbox/usb"
)

func errpanic(err error) {
	if err != nil {
		log.Fatalf("%+v", err)
	}
}

func main() {
	deviceInfo := func() usb.DeviceInfo {
		infos, err := usb.EnumerateHid(0, 0)
		errpanic(err)
		for _, di := range infos {
			if di.Serial == "" || di.Product == "" {
				continue
			}
			if di.Product == "bb02-bootloader" {
				return di
			}
		}
		panic("could no find a bitbox02")

	}()
	fmt.Println(deviceInfo)

	hidDevice, err := deviceInfo.Open()
	errpanic(err)

	const bitbox02BootloaderCMD = 0x80 + 0x40 + 0x03
	comm := u2fhid.NewCommunication(hidDevice, bitbox02BootloaderCMD)
	i := 0
	for {
		fmt.Println("query", i)
		response, err := comm.Query([]byte("v"))
		errpanic(err)
		fmt.Println(response)
		if !bytes.Equal(response, []byte{118, 0, 29, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0}) {
			panic("not equal")
		}
		i++
	}
}
