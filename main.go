package main

import (
	"LocalIPv4UseIPv6/model"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"github.com/songgao/packets/ethernet"
	"github.com/songgao/water"
	"io"
	"log"
	"net"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
)

func main() {
	//check runtime
	if os.Getuid() != 0 {
		log.Printf("no root mode! exit ...\n")
		return
	}

	argv := os.Args
	if len(argv) < 2 {
		log.Printf("need config file, exit ...\n")
		return
	}

	configFile, err := os.Open(argv[1])
	if err != nil {
		log.Printf("open config file %s, err: %s, exit ...\n", argv[1], err)
		return
	}

	// read config to str
	jsonData, err := io.ReadAll(configFile)
	if err != nil {
		log.Printf("read data err: %s\n", err)
	}
	_ = configFile.Close()

	println("config data = ", string(jsonData))

	// parse json
	var config model.Config
	_ = json.Unmarshal(jsonData, &config)

	// init tap device
	const TapName = "ipv6-ipv4"
	tapConfig := water.Config{
		DeviceType: water.TAP,
	}
	tapConfig.Name = TapName
	ifce, err := water.New(tapConfig)
	if err != nil {
		log.Printf("create tap device error: %s\n", err)
		return
	}

	// enable tap device
	cmd1 := exec.Command("ip", "addr", "add", fmt.Sprintf("%s/%s", config.Local.IP, config.Local.Mask), "dev", TapName)
	cmd2 := exec.Command("ip", "link", "set", "dev", TapName, "up")
	err = cmd1.Run()
	if err != nil {
		log.Printf("add ip to device err: %s\n", err)
		return
	}
	err = cmd2.Run()
	if err != nil {
		log.Printf("enable device err: %s\n", err)
		return
	}

	// create tcp server, to listen remote connect
	go socketListen(ifce)

	// accept data from user program and send data to remote
	go sendData(ifce, config)

	sig := make(chan os.Signal, 3)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGABRT, syscall.SIGHUP)
	<-sig

}

func socketListen(ifce *water.Interface) {
	listen, err := net.Listen("tcp6", ":1113")
	if err != nil {
		log.Printf("TCP Server create error %s, exit ...", err)
		os.Exit(-1)
	}
	for {
		conn, err := listen.Accept()
		if err != nil {
			log.Printf("accept error %s", err)
			continue
		}
		go processTCPClient(ifce, conn)
	}
}

func processTCPClient(ifce *water.Interface, conn net.Conn) {
	defer conn.Close()
	dataSize := make([]byte, 4)
	conn.Read(dataSize)
	log.Printf("read data from remote. size = %d\n", binary.LittleEndian.Uint32(dataSize))
	frameData := make([]byte, binary.LittleEndian.Uint32(dataSize))
	conn.Read(frameData)
	ifce.Write(frameData)
	log.Printf("write data to tap finish\n")
}

func sendData(ifce *water.Interface, config model.Config) {
	for {
		var frame = ethernet.Frame{}
		frame.Resize(1500)
		n, err := ifce.Read(frame)
		if err != nil {
			log.Printf("read data err: s%\n", err)
		}
		// one frame from ethernet
		frame = frame[:n]
		// get frame type ref: https://www.cnblogs.com/code1992/p/9829198.html
		frameType := frame.Ethertype()
		// get payload
		payload := frame.Payload()
		var targetIp []byte = nil
		if frameType[0] == 0x08 && frameType[1] == 0x06 {
			// arp
			targetIp = payload[24:28]
		}

		if frameType[0] == 0x08 && frameType[1] == 0x00 {
			// ip
			targetIp = payload[16:20]
		}

		if targetIp != nil {
			var targetIpv6 = ""
			ipStr := fmt.Sprintf("%d.%d.%d.%d", targetIp[0], targetIp[1], targetIp[2], targetIp[3])
			for _, item := range config.Remote {
				if item.Ipv4 == ipStr {
					targetIpv6 = item.Ipv6
					break
				}
			}
			if targetIpv6 != "" {
				go sendDataToRemote(frame, targetIpv6)
			} else {
				// not found target ipv6
				log.Printf("cannot found ipv6 address by ipv4(%s)\n", ipStr)
			}
		} else {
			// un support type
			log.Printf("un support type %x %x\n", frameType[0], frameType[1])
		}
	}
}

func sendDataToRemote(data []byte, address string) {
	log.Printf("send data to remote size: %d, address: %s", len(data), address)
	dataSize := make([]byte, 4)
	binary.LittleEndian.PutUint32(dataSize, uint32(len(data)))
	conn, err := net.Dial("tcp6", fmt.Sprintf("[%s]:1113", address))
	if err != nil {
		log.Printf("conn server failed, err:%s\n", err)
		return
	}
	defer conn.Close()
	targetData := append(dataSize, data...)
	targetSize := len(dataSize)
	log.Printf("send data to remote size: %d\n", targetSize)
	for targetSize > 0 {
		n, err := conn.Write(targetData)
		if err != nil && err != io.EOF {
			log.Printf("send data to target error: %s\n", err)
			break
		}
		log.Printf("send data size: %d\n", n)
		targetSize -= n
		if targetSize <= 0 || err == io.EOF {
			break
		}
		targetData = targetData[n:]
	}
}
