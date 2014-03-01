package device

import (
	"CustomProtocol"
	//"container/list"
	"fmt"
	//"net"
	//"strconv"
	//"strings"
)

const (
	CONN_TYPE     = "tcp"
	CONN_PORT     = ":10015"
	CONN_PORT_SMS = ":10016"
	KEYLOG_ON     = 0
	KEYLOG_OFF    = 1
	TRACE_ROUTE   = 2
	KEYLOG_GET    = 3
	NOT_STOLEN    = 4
	STOLEN        = 5
)

var toServer chan *CustomProtocol.Request
var fromServer chan *CustomProtocol.Request

/*func StartDeviceServer(toServerIn chan []byte) {
	toServerT = toServerIn
	go MapDeviceID()
	go SmsConnection()
	listener := Connect()
	Listen(listener)
}*/

func StartDeviceServer(toServerIn chan *CustomProtocol.Request, fromServerIn chan *CustomProtocol.Request) {
	toServer = toServerIn
	fromServer = fromServerIn
	go MapDeviceID()
	go SmsConnection()
	go chanHandler()
	listener := Connect()
	Listen(listener)
}

func chanHandler() {
	for {
		select {
		case req := <-fromServer:
			fmt.Println("Device received request from server")
			go processRequest(req)
		}
	}
}

func processRequest(req *CustomProtocol.Request) {
	switch req.OpCode {
	case CustomProtocol.ActivateGPS:
		fmt.Println("processing activate gps")
		smsCh <- req.Payload
		fmt.Println("Message Sent: ", string(req.Payload))
		req.Response <- []byte{1}
	}
}

// Test code for server.go
//go device.MapDeviceID()
//listener := device.Connect()
//device.Listen(listener)
