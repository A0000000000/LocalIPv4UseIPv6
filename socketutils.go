package main

import (
	"encoding/binary"
	"fmt"
	"github.com/songgao/water"
	"io"
	"log"
	"net"
	"os"
)

func revDataFromRemote(ifce *water.Interface) {
	listen, err := net.Listen("tcp6", ":1113")
	if err != nil {
		log.Printf("TCP Server create error %s, exit ...", err)
		os.Exit(-1)
	}
	for {
		conn, err := listen.Accept()
		if err != nil {
			continue
		}
		go progressConnection(ifce, conn)
	}
}

func progressConnection(ifce *water.Interface, conn net.Conn) {
	defer conn.Close()
	dataSize := make([]byte, 4)
	conn.Read(dataSize)
	log.Printf("read data size = %d, prepare send to tap device\n", binary.LittleEndian.Uint32(dataSize))
	frameData := make([]byte, binary.LittleEndian.Uint32(dataSize))
	conn.Read(frameData)
	writeData(ifce, frameData)
}

func sendData(targetIp []byte, data []byte, config Config) {
	if targetIp == nil {
		log.Println("target ip is nil")
		return
	}
	ipv6 := getIpv6FromIpv4(targetIp, config)
	if ipv6 == "" {
		log.Printf("cannot found ipv6 address by ipv4\n")
		return
	}
	sendDataToRemote(data, ipv6)
}

func sendDataToRemote(data []byte, address string) {
	dataSize := make([]byte, 4)
	binary.LittleEndian.PutUint32(dataSize, uint32(len(data)))
	conn, err := net.Dial("tcp6", fmt.Sprintf("[%s]:1113", address))
	if err != nil {
		log.Printf("conn server failed, err:%s\n", err)
		return
	}
	defer conn.Close()
	targetData := append(dataSize, data...)
	targetSize := len(targetData)
	log.Printf("prepare send data to remote, size: %d\n", targetSize)
	for targetSize > 0 {
		n, err := conn.Write(targetData)
		if err != nil && err != io.EOF {
			log.Printf("send data to target error: %s\n", err)
			break
		}
		targetSize -= n
		if targetSize <= 0 || err == io.EOF {
			break
		}
		targetData = targetData[n:]
	}
}
