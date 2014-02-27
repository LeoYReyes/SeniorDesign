/*
 * @author Nathan Plotts (nwp0002@auburn.edu)
 * This file is where the "start" function is for both the LaptopHub and
 * the GPSHub connection handling.
 */

package device

import (
	"CustomRequest"
	//"container/list"
	//"fmt"
	//"net"
	//"strconv"
	//"strings"
)

/* Constants for connections and OP codes are stored here, accessible by
 * all in the device package.
 */
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

// Global channel variables for sending and receiving request to/from the server
var toServer chan *CustomRequest.Request
var fromServer chan *CustomRequest.Request

/*
 * This method is the function that you call to start both the laptop and GPS connection
 * handling. Once connection handling has started then connections will be continually
 * accepted and requests will continually be handled.
 */
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
