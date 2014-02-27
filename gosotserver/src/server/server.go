package main

import (
	"CustomRequest"
	"device"
	//"fmt"
	//"time"
	"webserver"
)

var toWebCh = make(chan *CustomRequest.Request)
var fromWebCh = make(chan *CustomRequest.Request)
var toDatabaseCh = make(chan *CustomRequest.Request)
var fromDatabaseCh = make(chan *CustomRequest.Request)
var toDeviceCh = make(chan *CustomRequest.Request)
var fromDeviceCh = make(chan *CustomRequest.Request)

var testDeviceCh = make(chan []byte)
var testWebCh = make(chan []byte)

func main() {
	// channel can take optional capacity param to make it asynchronous
	//comChannel := make(chan string)

	//go webserver.StartWebServer(fromWebCh, toWebCh)
	//go device.StartDeviceServer(fromDeviceCh, toDeviceCh)

	// FOR TESTING
	go webserver.StartWebServer(testWebCh)
	go device.StartDeviceServer(testDeviceCh)
	for {
		select {
		case c := <-testDeviceCh:
			testWebCh <- c
		}
	}
	//fmt.Println(<-fromWebCh)
	//time.Sleep(3000000 * time.Millisecond)
	//TODO: figure out a way to leave it running with for loop
}
