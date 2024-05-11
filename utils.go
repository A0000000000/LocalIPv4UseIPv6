package main

import (
	"encoding/json"
	"fmt"
	"github.com/songgao/water"
	"io"
	"log"
	"os"
	"os/exec"
)

func checkProgramRuntime() (bool, string) {
	if os.Getuid() != 0 {
		return false, "not in root mode"
	}

	argv := os.Args
	if len(argv) < 2 {
		return false, "no config file"
	}
	return true, ""
}

func getConfig() (error, Config) {
	var config Config
	argv := os.Args
	configFile, err := os.Open(argv[1])
	if err != nil {
		return err, config
	}
	jsonData, err := io.ReadAll(configFile)
	if err != nil {
		return err, config
	}
	_ = configFile.Close()
	_ = json.Unmarshal(jsonData, &config)
	return nil, config
}

func createAndInitTapDevice(config Config) (error, *water.Interface) {
	const TapName = "ipv6-ipv4"
	tapConfig := water.Config{
		DeviceType: water.TAP,
	}
	tapConfig.Name = TapName
	ifce, err := water.New(tapConfig)
	if err != nil {
		return err, nil
	}
	cmd1 := exec.Command("ip", "addr", "add", fmt.Sprintf("%s/%s", config.Local.IP, config.Local.Mask), "dev", TapName)
	cmd2 := exec.Command("ip", "link", "set", "dev", TapName, "up")
	err = cmd1.Run()
	if err != nil {
		return err, nil
	}
	err = cmd2.Run()
	if err != nil {
		log.Printf("enable device err: %s\n", err)
		return err, nil
	}
	return nil, ifce
}

func getIpv6FromIpv4(ipv4 []byte, config Config) string {
	var ipv6 = ""
	ipStr := fmt.Sprintf("%d.%d.%d.%d", ipv4[0], ipv4[1], ipv4[2], ipv4[3])
	for _, item := range config.Remote {
		if item.Ipv4 == ipStr {
			ipv6 = item.Ipv6
			break
		}
	}
	return ipv6
}
