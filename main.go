package main

import (
	"andtest/logger"
	"andtest/scrcpy"
)

func main() {
	tag := "webrtc-proxy"
	logger.LogI(tag, "Service Started")
	//
	////for {
	////	network.Network(tag)
	////	time.Sleep(20 * time.Second)
	////}
	//socket.Socket(tag)
	//select {}

	//control.Control(tag)

	scrcpy.Scrcpy(tag)
}
