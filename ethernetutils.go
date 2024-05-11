package main

import "github.com/songgao/packets/ethernet"

func parseFrameData(frame ethernet.Frame) ([2]byte, []byte) {
	frameType := frame.Ethertype()
	payload := frame.Payload()
	var targetIp []byte

	// get frame type ref: https://www.cnblogs.com/code1992/p/9829198.html
	if frameType[0] == 0x08 && frameType[1] == 0x06 {
		// arp
		targetIp = payload[24:28]
	}
	if frameType[0] == 0x08 && frameType[1] == 0x00 {
		// ip
		targetIp = payload[16:20]
	}
	// Todo: other agreement
	return frameType, targetIp
}
