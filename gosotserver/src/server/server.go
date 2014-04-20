package main

import (
	"CustomProtocol"
	"databaseSOT"
	"device"
	//"fmt"
	//"time"
	"webserver"
)

var toWebCh = make(chan *CustomProtocol.Request, 1000)
var fromWebCh = make(chan *CustomProtocol.Request, 1000)
var toDatabaseCh = make(chan *CustomProtocol.Request, 1000)
var fromDatabaseCh = make(chan *CustomProtocol.Request, 1000)
var toDeviceCh = make(chan *CustomProtocol.Request, 1000)
var fromDeviceCh = make(chan *CustomProtocol.Request, 1000)

var testDeviceCh = make(chan []byte)
var testWebCh = make(chan []byte)

func main() {
	// channel can take optional capacity param to make it asynchronous
	//comChannel := make(chan string)

	go webserver.StartWebServer(fromWebCh, toWebCh)
	go device.StartDeviceServer(fromDeviceCh, toDeviceCh)
	go databaseSOT.StartDatabaseServer(fromDatabaseCh, toDatabaseCh)

	for {
		select {
		case req := <-fromWebCh:
			reRoute(req)
		case req := <-fromDeviceCh:
			reRoute(req)
		case req := <-fromDatabaseCh:
			reRoute(req)
		}

	}
	//TODO: figure out a way to leave it running with for loop
}

func reRoute(req *CustomProtocol.Request) {
	switch req.Destination {
	case CustomProtocol.Database:
		toDatabaseCh <- req
		//fmt.Println("Reroute to database")
	case CustomProtocol.Web:
		toWebCh <- req
		//fmt.Println("Reroute to web")
	case CustomProtocol.DeviceGPS, CustomProtocol.DeviceLaptop:
		toDeviceCh <- req
		//fmt.Println("Reroute to device")
	}
}
