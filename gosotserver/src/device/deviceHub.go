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
	"strings"
)

/*
 * All constants for the entire device package are defined here.
 */
const (
	CONN_TYPE     = "tcp"
	CONN_PORT     = ":10015"
	CONN_PORT_SMS = ":10016"
	//KEYLOG_ON     = 0
	//KEYLOG_OFF    = 1
	//TRACE_ROUTE   = 2
	//KEYLOG_GET    = 3
	//NOT_STOLEN    = 4
	//STOLEN        = 5
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
	//go MapDeviceID()
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
			go processRequest(req) //todo go routine here may not be necessary
		}
	}
}

/*
 * Remove '[', ']', and '|' from params
 */
func sanitizeGPSInput(params []string) []string { //todo
	for i := 0; i < len(params); i++ {
		params[i] = strings.Replace(params[i], "[", "", -1)
		params[i] = strings.Replace(params[i], "]", "", -1)
		params[i] = strings.Replace(params[i], "|", "", -1)
	}
	return params
}

/*
 * This method takes requests from the fromServer request channel and parses the request.
 * It then uses the OpCode from the request to reroute the request to the correct hub.
 *
 * GPS req payload structure (esc character delimiter):
 * <phone number><PIN><variable numnber of params>
 * exception: FreestyleMsg does not require a PIN, as it lets the message
 * be completely customizable.
 * note: '[', ']', and '|' are reserved as delimiters and should not be included
 * in any parameters
 */
func processRequest(req *CustomProtocol.Request) {
	switch req.OpCode {
	//params: phone number, PIN
	case CustomProtocol.ActivateGPS:
		fmt.Println("processing activate gps")
		payload := sanitizeGPSInput(CustomProtocol.ParsePayload(req.Payload))
		if len(payload) >= 2 {
			msg := "[" + payload[0] + "]" + payload[1] + ".0.|"
			smsCh <- msg
			fmt.Println("Message Sent: ", msg)
			req.Response <- []byte{1}
		} else {
			req.Response <- []byte{0}
		}
	//activates geofence 1
	//params: phone number, PIN, active/deactive (1/0), radius (feet)
	case CustomProtocol.ActivateGeofence:
		fmt.Println("processing activate geofence")
		payload := sanitizeGPSInput(CustomProtocol.ParsePayload(req.Payload))
		if len(payload) >= 4 {
			msg := "[" + payload[0] + "]" + payload[1] + ".2." + payload[2] + ".1.0." + payload[3] + ".|"
			smsCh <- msg
			fmt.Println("Message Sent: ", msg)
			req.Response <- []byte{1}
		} else {
			req.Response <- []byte{0}
		}
	//params: phone number, PIN
	case CustomProtocol.SleepGeogram:
		fmt.Println("processing sleep geogram")
		payload := sanitizeGPSInput(CustomProtocol.ParsePayload(req.Payload))
		if len(payload) >= 2 {
			//msg := "[" + payload[0] + "]" + payload[1] + ".1.|"
			//new sleeping method, awake for 120 seconds, asleep for 86400 (24 hours)
			//also wakes on SMS or motion
			msg := "[" + payload[0] + "]" + payload[1] + ".5.120.86400.6.|"
			smsCh <- msg
			fmt.Println("Message Sent: ", msg)
			req.Response <- []byte{1}
		} else {
			req.Response <- []byte{0}
		}
	//params: phone number, PIN, interval (seconds) (0 to disable)
	case CustomProtocol.ActivateIntervalGps:
		fmt.Println("processing activate interval gps")
		payload := sanitizeGPSInput(CustomProtocol.ParsePayload(req.Payload))
		if len(payload) >= 3 {
			msg := "[" + payload[0] + "]" + payload[1] + ".4." + payload[2] + ".|"
			smsCh <- msg
			fmt.Println("Message Sent: ", msg)
			req.Response <- []byte{1}
		} else {
			req.Response <- []byte{0}
		}
	// sets location for geofence 1
	// params: phone number, PIN, latitude format ddmm.mmmm without the '.',
	// longitude format dddmm.mmmm without the '.'
	case CustomProtocol.SetGeofence:
		fmt.Println("processing set geofence")
		payload := sanitizeGPSInput(CustomProtocol.ParsePayload(req.Payload))
		if len(payload) >= 4 {
			//lat
			latMsg := "[" + payload[0] + "]" + payload[1] + ".6.128." + payload[2] + ".|"
			smsCh <- latMsg
			fmt.Println("Message Sent: ", latMsg)
			//long
			longMsg := "[" + payload[0] + "]" + payload[1] + ".6.132." + payload[3] + ".|"
			smsCh <- longMsg
			fmt.Println("Message Sent: ", longMsg)
			req.Response <- []byte{1}
		} else {
			req.Response <- []byte{0}
		}
	//params: phone number, PIN
	case CustomProtocol.GeogramSetup:
		fmt.Println("processing geogram setup")
		payload := sanitizeGPSInput(CustomProtocol.ParsePayload(req.Payload))
		if len(payload) >= 2 {
			//motion alert
			motMsg := "[" + payload[0] + "]" + payload[1] + ".6.200." + MOTION_ALERT + ".|"
			smsCh <- motMsg
			fmt.Println("Message Sent: ", motMsg)

			//hyperlink1
			hyp1Msg := "[" + payload[0] + "]" + payload[1] + ".6.450." + HYPERLINK_1 + ".|"
			smsCh <- hyp1Msg
			fmt.Println("Message Sent: ", hyp1Msg)

			//hyperlink2
			hyp2Msg := "[" + payload[0] + "]" + payload[1] + ".6.500." + HYPERLINK_2 + ".|"
			smsCh <- hyp2Msg
			fmt.Println("Message Sent: ", hyp2Msg)

			//hyperlink3
			hyp3Msg := "[" + payload[0] + "]" + payload[1] + ".6.550." + HYPERLINK_3 + ".|"
			smsCh <- hyp3Msg
			fmt.Println("Message Sent: ", hyp3Msg)

			//geofence alert
			geoMsg := "[" + payload[0] + "]" + payload[1] + ".6.250." + GEOFENCE_ALERT + ".|"
			smsCh <- geoMsg
			fmt.Println("Message Sent: ", geoMsg)

			req.Response <- []byte{1}
		} else {
			req.Response <- []byte{0}
		}
	//params: phone number, message
	case CustomProtocol.FreestyleMsg:
		fmt.Println("processing freestyle msg")
		payload := sanitizeGPSInput(CustomProtocol.ParsePayload(req.Payload))
		if len(payload) >= 2 {
			msg := "[" + payload[0] + "]" + payload[1] + "|"
			smsCh <- msg
			fmt.Println("Message Sent: ", msg)
			req.Response <- []byte{1}
		} else {
			req.Response <- []byte{0}
		}
	case CustomProtocol.UpdateUserKeylogData: //todo idk if it is listening for a response. also do the NO_OP messages interefere with getmessage?
		go ProcessLapReq(req) //todo is creating a thread for this a good idea?
	case CustomProtocol.UpdateUserIPTraceData: //todo include IP of connection in the trace route, currently starts at first node from router
		go ProcessLapReq(req)
	case CustomProtocol.FlagStolen:
		go ProcessLapReq(req)
	case CustomProtocol.FlagNotStolen:
		go ProcessLapReq(req)
	default:
		req.Response <- []byte{0}
	}
}
