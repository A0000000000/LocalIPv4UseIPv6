package main

import (
	"github.com/songgao/packets/ethernet"
	"github.com/songgao/water"
	"log"
)

func writeData(ifce *water.Interface, data []byte) {
	ifce.Write(data)
}

func readData(ifce *water.Interface, config Config) {
	for {
		var frame = ethernet.Frame{}
		frame.Resize(1500)
		n, err := ifce.Read(frame)
		if err != nil {
			log.Printf("read data err: %s\n", err)
		}
		frame = frame[:n]
		_, targetIp := parseFrameData(frame)
		go sendData(targetIp, frame, config)
	}
}
