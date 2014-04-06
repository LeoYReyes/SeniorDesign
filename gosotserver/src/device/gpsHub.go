/*
 * @author Nathan Plotts (nwp0002@auburn.edu)
 * @author Charlie Baker (cmb0049@auburn.edu)
 * @author Leo Reyes (lyr0001@auburn.edu)
 * This file is where the the gpsHub struct is defined and also
 * where general Geogram connection handling functions will be
 * stored.
 */

package device

import (
	"CustomProtocol"
	//"container/list"
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"
)

const (
	INTERVAL_TIME     = "60"
	MOTION_AWAKE_TIME = 300000 * time.Millisecond
)

var smsConn net.Conn
var smsCh = make(chan string)

/*
 * This method creates a connection which creates a new socket, opens the port
 * that connections go through, and returns a listener.
 */
func GPSConnect() net.Listener {
	//connect
	listener, err := net.Listen(CONN_TYPE, CONN_PORT_SMS)
	if err != nil {
		fmt.Println("Error listening:", err.Error())
	} else {

	}
	fmt.Println("Connection created on " + CONN_TYPE + " " + CONN_PORT_SMS)
	return listener
}

/*
 * This method takes in the listener object created by the GPSConnect function
 * and begins accepting connections through it. After creating a connection
 * with a device it then calls the GPSCommunicate function.
 */
func GPSListen(listener net.Listener) {
	for {
		smsConn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error connecting", err)
		}
		fmt.Println("Connection established with SMS client")
		// make this a go routine and it likely already
		// has the functionality to suppport multiple phones
		GPSCommunicate(smsConn)
	}
}

/*
 * This method takes in the connection created in the GPSListen function and
 * receives GPS coordinates through it. Any coordinates it receives it then
 * sends to the database.
 */
func GPSCommunicate(conn net.Conn) {
	buffer := make([]byte, 512)
	msg := ""
	lastPing := (time.Now()).Unix()
	for {
		remaining := ""
		select {
		case m := <-smsCh:
			fmt.Println("smsCh: " + m)
			conn.Write([]byte(m))
			break
		default:
			//fmt.Println("Waiting to read from smsdevice")
			bytesRead, _ := conn.Read(buffer)
			if bytesRead > 0 {
				if buffer[0] == '|' {
					//fmt.Println("Heartbeat <3")
					conn.Write([]byte("|")) //heartbeat response to ensure connection is alive
					lastPing = (time.Now()).Unix()
				} else {
					received := remaining + string(buffer[0:bytesRead])
					index := strings.Index(received, "|")
					if index == -1 {
						remaining = received
					} else if index+1 < len(received) {
						remaining = received[index:]
					} else {
						remaining = ""
					}
					fmt.Println("Received msg: ", received)
					number := parsePhoneNumber(received)
					msg = googleMapLinkParser(received)
					//react based on message
					if msg != "" { //try to parse it to coords first, if it fails it is another type of message
						fmt.Println("parsed msg: ", msg)
						msg = strings.Replace(msg, ",", string(0x1B), -1)
						msg = number + string(0x1B) + msg + string(0x1B)
						go UpdateMapCoords(msg)
					} else if strings.Contains(received, MOTION_ALERT) {
						motionAlert(number)
					} else if strings.Contains(received, GEOFENCE_ALERT) {
						//geofenceAlert(number) //todo uncomment when functional
					} else {
						fmt.Println("Message format not recognized")
					}
				}
			}
		}
		if (time.Now()).Unix()-lastPing > 30 {
			fmt.Println("Lost connection with SMS client")
			conn.Close()
			break
		}
	}
}

func parsePhoneNumber(msg string) string {
	indexStart := strings.Index(msg, "[")
	indexEnd := strings.Index(msg, "]")
	if indexStart > -1 && indexEnd > -1 && indexEnd > indexStart {
		return msg[indexStart+1 : indexEnd]
	}
	return ""
}

func motionAlert(phoneNumber string) {
	fmt.Println(phoneNumber + " " + MOTION_ALERT)
	go motionAlertTimer(phoneNumber)
	/*
		//report stolen
		payload := append([]byte(phoneNumber), 0x1B)
		response := make(chan []byte)
		req := &CustomProtocol.Request{Id: CustomProtocol.AssignRequestId(), Destination: CustomProtocol.Database,
			Source: CustomProtocol.DeviceGPS, OpCode: CustomProtocol.ActivateGPS, Payload: payload,
			Response: response}
		toServer <- req
		//add response check later
		//interval gps request
		pin := getPIN(phoneNumber)
		interval := INTERVAL_TIME

		payload2 := []byte(phoneNumber)
		payload2 = append(payload2, 0x1B)
		payload2 = append(payload2, []byte(pin)...)
		payload2 = append(payload2, 0x1B)
		payload2 = append(payload2, []byte(interval)...)
		payload2 = append(payload2, 0x1B)
		response2 := make(chan []byte)
		req2 := &CustomProtocol.Request{Id: CustomProtocol.AssignRequestId(), Destination: CustomProtocol.DeviceGPS,
			Source: CustomProtocol.DeviceGPS, OpCode: CustomProtocol.ActivateIntervalGps, Payload: payload2,
			Response: response2}
		toServer <- req2
		//add response check later
	*/
}

func motionAlertTimer(phoneNumber string) {
	time.Sleep(MOTION_AWAKE_TIME)

	//check if stolen
	buf := []byte(phoneNumber)
	buf = append(buf, 0x1B)

	response := make(chan []byte)
	req := &CustomProtocol.Request{Id: CustomProtocol.AssignRequestId(), Destination: CustomProtocol.Database,
		Source: CustomProtocol.DeviceGPS, OpCode: CustomProtocol.CheckDeviceStolen, Payload: buf,
		Response: response}
	toServer <- req
	isStolen := <-response

	if isStolen[0] != 1 {
		//send sleep command
		pin := getPIN(phoneNumber)

		payload := []byte(phoneNumber)
		payload = append(payload, 0x1B)
		payload = append(payload, []byte(pin)...)
		payload = append(payload, 0x1B)

		response2 := make(chan []byte)
		req2 := &CustomProtocol.Request{Id: CustomProtocol.AssignRequestId(), Destination: CustomProtocol.DeviceGPS,
			Source: CustomProtocol.DeviceGPS, OpCode: CustomProtocol.SleepGeogram, Payload: payload,
			Response: response2}
		toServer <- req2
		fmt.Println(phoneNumber + " has not left geofence. Sleeping")
	}
}

func geofenceAlert(phoneNumber string) { //todo add functionality
	fmt.Println(phoneNumber + " " + MOTION_ALERT)
	//report stolen
	payload := append([]byte(phoneNumber), 0x1B)
	response := make(chan []byte)
	req := &CustomProtocol.Request{Id: CustomProtocol.AssignRequestId(), Destination: CustomProtocol.Database,
		Source: CustomProtocol.DeviceGPS, OpCode: CustomProtocol.ActivateGPS, Payload: payload,
		Response: response}
	toServer <- req
	//add response check later
	//interval gps request
	pin := getPIN(phoneNumber)
	interval := INTERVAL_TIME

	payload2 := []byte(phoneNumber)
	payload2 = append(payload2, 0x1B)
	payload2 = append(payload2, []byte(pin)...)
	payload2 = append(payload2, 0x1B)
	payload2 = append(payload2, []byte(interval)...)
	payload2 = append(payload2, 0x1B)
	response2 := make(chan []byte)
	req2 := &CustomProtocol.Request{Id: CustomProtocol.AssignRequestId(), Destination: CustomProtocol.DeviceGPS,
		Source: CustomProtocol.DeviceGPS, OpCode: CustomProtocol.ActivateIntervalGps, Payload: payload2,
		Response: response2}
	toServer <- req2
	//add response check later
}

func getPIN(phoneNumber string) string {
	return "1234"
}

/*
 * This method is the "main" method for the gpsHub file. When it's called it
 * begins and keeps GPS tracker communications open and running
 */
func SmsConnection() {
	//send & receive
	listener := GPSConnect()
	GPSListen(listener)
}

/*
 * This method takes in string that was received in a text message and parses
 * out the GPS coordinates, and then returns the coordinates as a string.
 */
func googleMapLinkParser(input string) string {
	result := ""
	str := input
	//str := "[1111111111]http://maps.google.com/maps?q=32+36.3143,-085+29.1954+()&z=19|"
	index := strings.Index(str, "=")
	if index == -1 {
		return ""
	}
	str = str[index+1:]
	// latitude
	indexComma := strings.Index(str, ",")
	if indexComma == -1 {
		return ""
	}
	lat := str[0:indexComma]

	indexPlus := strings.Index(lat, "+")
	if indexPlus == -1 || indexPlus > indexComma {
		return ""
	}
	latWhole := lat[0:indexPlus]

	latDecimal, err1 := strconv.ParseFloat(lat[indexPlus+1:indexComma], 16)
	if err1 != nil {
		return ""
	}

	//longitude
	long := str[indexComma+1:] //indexComma is strill from the lat calculationsabove

	indexEnd := strings.Index(long, "+(")
	if indexEnd == -1 {
		return ""
	}
	lat = str[0:indexEnd]
	indexPlus = strings.Index(long, "+")
	if indexPlus == -1 || indexPlus > indexEnd {
		return ""
	}
	longWhole := long[0:indexPlus]
	longDecimal, err2 := strconv.ParseFloat(long[indexPlus+1:indexEnd], 16)
	if err2 != nil {
		return ""
	}
	latDecimal = latDecimal / 60
	longDecimal = longDecimal / 60
	latStr := []byte{}
	longStr := []byte{}
	latStr = strconv.AppendFloat(latStr, latDecimal, 'g', 4, 32)
	longStr = strconv.AppendFloat(longStr, longDecimal, 'g', 4, 32)
	if len(latStr) < 3 || len(longStr) < 3 {
		return ""
	}
	result = strings.Join([]string{latWhole, ".", string(latStr[2:]), ",", longWhole, ".", string(longStr[2:])}, "")
	return result
}
