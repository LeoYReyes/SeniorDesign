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

var toServerT chan []byte

func GPSConnect() net.Listener {
	listener, err := net.Listen(CONN_TYPE, CONN_PORT_SMS)
	if err != nil {
		fmt.Println("Error listening:", err.Error())
	} else {

	}
	fmt.Println("Connection created on " + CONN_TYPE + " " + CONN_PORT_SMS)
	return listener
}

func GPSListen(listener net.Listener) net.Conn {
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error connecting", err)
		}
		fmt.Println("Connection established with SMS client")
		GPSGetMessages(conn)
	}
}

func GPSGetMessages(conn net.Conn) {
	msg := ""
	for {
		buffer := make([]byte, 512)
		bytesRead, _ := conn.Read(buffer)
		if bytesRead > 0 {
			if bytesRead > 10 {
				received := string(buffer[0:bytesRead])
				msg = googleMapLinkParser(received)
				fmt.Println("Received msg: ", msg)
				req := &CustomProtocol.Request{Id: CustomProtocol.AssignRequestId(), Destination: CustomProtocol.Web, Source: CustomProtocol.DeviceGPS,
					OpCode: CustomProtocol.UpdateWebMap, Payload: []byte(msg)}
				toServer <- req
			} else {
				conn.Write([]byte("|"))
			}
		}
	}
}

func SmsConnection() {
	//connect
	listener := GPSConnect()
	GPSListen(listener)
	//send & receive

}

func googleMapLinkParser(input string) string {
	result := ""
	str := input
	//str := "[1111111111]http://maps.google.com/maps?q=32+36.3143,-085+29.1954+()&z=19|"
	str = str[strings.Index(str, "=")+1:]
	latDecimal, _ := strconv.ParseFloat(str[3:10], 16)
	longDecimal, _ := strconv.ParseFloat(str[16:23], 16)
	latDecimal = latDecimal / 60
	longDecimal = longDecimal / 60
	latStr := []byte{}
	longStr := []byte{}
	latStr = strconv.AppendFloat(latStr, latDecimal, 'g', 4, 32)
	longStr = strconv.AppendFloat(longStr, longDecimal, 'g', 4, 32)
	result = strings.Join([]string{result, str[0:2], ".", string(latStr[2:]), str[10:15], ".", string(longStr[2:])}, "")
	return result
}
