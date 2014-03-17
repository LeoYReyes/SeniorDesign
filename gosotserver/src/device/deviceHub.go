/*
 * @author Nathan Plotts (nwp0002@auburn.edu)
 * This file is the main controller for both the laptopHub and gpsHub. When StartDeviceServer
 * is called both laptopHub and gpsHub create their sockets and begin accepting connections.
 * The channels for sending to the server are initialized here. The channels carry byte arrays
 * which hold our requests for both the database and the webClient. Request structure is defined
 * in Request.go.
 *
 * We also store all the constant values for both hubs in this file. This allows us to keep
 * them synced across all hubs as some constants are used in both places.
 */

package device

import (
	"CustomProtocol"
	"fmt"
)

/*
 * All constants for the entire device package are defined here.
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

/*
 * These channels are used for sending and receiving requests to and from the server.
 */
var toServer chan *CustomProtocol.Request
var fromServer chan *CustomProtocol.Request

/*
 * This is the main function for deviceHub. When this function is called it starts the
 * connection handlers for both the laptopHub and gpsHub in their own threads. Also,
 * this function starts a thread that is used to hash laptopDevice connections by their
 * ID so they can later be accessed.
 *
 * This method is called with 2 channel parameters. These channels come from, and are
 * defined, in the main function in server.go. These channels are then mapped to the
 * toServer and fromServer channel variables which are used to send and receive requests
 * to and from the rest of the server components.
 */
func StartDeviceServer(toServerIn chan *CustomProtocol.Request, fromServerIn chan *CustomProtocol.Request) {
	toServer = toServerIn
	fromServer = fromServerIn
	go MapDeviceID()
	go SmsConnection()
	go chanHandler()
	listener := Connect()
	Listen(listener)
}

/*
 * This method runs in its own thread and constantly reads from the fromServer channel.
 * As soon as it receives a request it calls the processRequest function and correctly
 * forwards and handles the request.
 */
func chanHandler() {
	for {
		select {
		case req := <-fromServer:
			fmt.Println("Device received request from server")
			go processRequest(req)
		}
	}
}

/*
 * This method takes requests from the fromServer request channel and parses the request.
 * It then uses the OpCode from the request to reroute the request to the correct hub.
 */
func processRequest(req *CustomProtocol.Request) {
	switch req.OpCode {
	case CustomProtocol.ActivateGPS:
		fmt.Println("processing activate gps")
		smsCh <- req.Payload
		fmt.Println("Message Sent: ", string(req.Payload))
		req.Response <- []byte{1}
	case CustomProtocol.UpdateUserKeylogData:
		go ProcessLapReq(req) //todo is creating a thread for this a good idea?
	case CustomProtocol.UpdateUserIPTraceData:
		go ProcessLapReq(req)
	}
}
