package device

import (
	"CustomRequest"
	//"container/list"
	//"fmt"
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

var toServer chan *CustomRequest.Request
var fromServer chan *CustomRequest.Request

func StartDeviceServer(toServerIn chan []byte) {
	toServerT = toServerIn
	go MapDeviceID()
	go SmsConnection()
	listener := Connect()
	Listen(listener)
}

/*func StartDeviceServer(fromServerIn chan *CustomRequest.Request, toServerIn chan *CustomRequest.Request) {
	toServer = toServerIn
	fromServer = fromServerIn
	go MapDeviceID()
	go SmsConnection()
	listener := Connect()
	Listen(listener)
}*/

// Test code for server.go
//go device.MapDeviceID()
//listener := device.Connect()
//device.Listen(listener)
