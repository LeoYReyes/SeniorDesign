package main

import (
	"CustomProtocol"
	"device"
	//"fmt"
	//"time"
	"webserver"
)

var toWebCh = make(chan *CustomProtocol.CustomProtocol)
var fromWebCh = make(chan *CustomProtocol.CustomProtocol)
var toDatabaseCh = make(chan *CustomProtocol.CustomProtocol)
var fromDatabaseCh = make(chan *CustomProtocol.CustomProtocol)
var toDeviceCh = make(chan *CustomProtocol.CustomProtocol)
var fromDeviceCh = make(chan *CustomProtocol.CustomProtocol)

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
		case req := <-testDeviceCh:
			testWebCh <- req
		case req := <-fromWebCh:
			reRoute(req)
		case req := <-fromDeviceCh:
			reRoute(req)
		case req := <-fromDatabaseCh:
			reRoute(req)
		}

	}
	//fmt.Println(<-fromWebCh)
	//time.Sleep(3000000 * time.Millisecond)
	//TODO: figure out a way to leave it running with for loop
}

func reRoute(req *CustomProtocol.CustomProtocol) {
	switch req.Destination {
	case CustomProtocol.Database:
		toDatabaseCh <- req
	case CustomProtocol.Web:
		toWebCh <- req
	case CustomProtocol.DeviceGPS, CustomProtocol.DeviceLaptop:
		toDeviceCh <- req
	}
}
