package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	res, reason := checkProgramRuntime()
	if !res {
		log.Println(reason)
		return
	}

	err, config := getConfig()
	if err != nil {
		log.Println(err)
		return
	}

	err, ifce := createAndInitTapDevice(config)
	if err != nil {
		log.Println(err)
		return
	}

	go revDataFromRemote(ifce)
	go readData(ifce, config)

	sig := make(chan os.Signal, 3)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGABRT, syscall.SIGHUP)
	<-sig
}
