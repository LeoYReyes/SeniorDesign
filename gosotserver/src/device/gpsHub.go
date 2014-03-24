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
	for {
		select {
		case m := <-smsCh:
			fmt.Println("smsCh: " + m)
			conn.Write([]byte(m))
		default:
			//fmt.Println("Waiting to read from smsdevice")
			bytesRead, _ := conn.Read(buffer)
			if bytesRead > 0 {
				if bytesRead > 10 {
					received := string(buffer[0:bytesRead])
					msg = googleMapLinkParser(received)
					fmt.Println("Received msg: ", msg)
					req := &CustomProtocol.Request{Id: CustomProtocol.AssignRequestId(), Destination: CustomProtocol.Web, Source: CustomProtocol.DeviceGPS,
						OpCode: CustomProtocol.UpdateWebMap, Payload: []byte(msg)}
					msg = strings.Replace(params[i], ",", 0x1B, -1)
					msg = msg + 0x1B
					//TODO send GPS coordinates to the database
					toServer <- req
					fmt.Println("Req sent to server")
				} else if buffer[0] == '|' {
					//fmt.Println("Heartbeat <3")
					conn.Write([]byte("|")) //heartbeat response to ensure connection is alive
				}
			}
		}
	}
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
	latDecimal, err1 := strconv.ParseFloat(str[3:10], 16)
	longDecimal, err2 := strconv.ParseFloat(str[16:23], 16)
	if err1 != nil || err2 != nil {
		return ""
	}
	latDecimal = latDecimal / 60
	longDecimal = longDecimal / 60
	latStr := []byte{}
	longStr := []byte{}
	latStr = strconv.AppendFloat(latStr, latDecimal, 'g', 4, 32)
	longStr = strconv.AppendFloat(longStr, longDecimal, 'g', 4, 32)
	result = strings.Join([]string{result, str[0:2], ".", string(latStr[2:]), str[10:15], ".", string(longStr[2:])}, "")
	return result
}
