/*
 * @author Nathan Plotts (nwp0002@auburn.edu)
 * This file is where the the laptopHub struct is defined and
 * also where general laptop connection handling functions will
 * be stored.
 */

package device

import (
	//"CustomProtocol"
	//"container/list"
	"fmt"
	"net"
	"strconv"
	//"strings"
)

type laptopHub struct {
	connections map[string]net.Conn
	//	DeviceBuffer
	mapDeviceQueue chan *deviceConnection
}

type deviceConnection struct {
	ld   LaptopDevice
	conn net.Conn
}

var lh = laptopHub{
	mapDeviceQueue: make(chan *deviceConnection, 20),
	connections:    make(map[string]net.Conn),
}

var deviceConn = new(deviceConnection)

/*
 * This method creates the socket that the server will be listening on for laptop
 * connections. It also opens the port on the server.
 */
func Connect() net.Listener {
	listener, err := net.Listen(CONN_TYPE, CONN_PORT)
	if err != nil {
		fmt.Println("Error listening:", err.Error())
		return listener
	}
	fmt.Println("Connection created on " + CONN_TYPE + " " + CONN_PORT)
	return listener
}

/*
 * This method begins accepting new connections from laptop devices. As connections
 * are opened they are handed off to the GetDeviceID function in a GoRoutine to be
 * read from. Reading the connection in a new thread negates worries of IO blocking.
 */
func Listen(listener net.Listener) {
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error connecting", err)
		}
		fmt.Println("Connection established with client")
		go GetDeviceID(conn)
	}
}

/*
 * This method is where long messages sent from a laptop device are read. They are
 * parsed as a string and then the OP code from the message is read. The message
 * handling is the handed to the correct handling function.
 */
func GetMessage(deviceConn *deviceConnection) {
	buffer := make([]byte, 10240)
	bytesRead, err := deviceConn.conn.Read(buffer)
	if err != nil {
		fmt.Println("Error reading", err)
	}
	msg := string(buffer[0:bytesRead])
	opCode, err := strconv.Atoi(msg[0:1])
	if err != nil {
		fmt.Println("Invalid OP code", err)
	} else {
		switch opCode {
		case TRACE_ROUTE:
			UpdateTraceroute(deviceConn, msg)
		case KEYLOG_GET:
			UpdateKeylog(deviceConn, msg)
		}
	}
}

/*
 * This method is called when a message's OP code is set to the TRACE_ROUTE
 * constant. It then takes the remaining string that consists of IP addresses
 * and parses them into a list.List object. The List is then added to the
 * client's list of TraceRoutes and a request is sent to the database to
 * sync the new list there.
 */
func UpdateTraceroute(deviceConn *deviceConnection, msg string) {
	/*newList := new(list.List)
	start := 1
	for i := 1; i < len(msg)-1; i++ {
		if msg[i:i+1] == "~" {
			newList.PushBack(msg[start:i])
			start = i + 1
		}
	}
	newList.PushBack(msg[start : len(msg)-1]) //to get the final IP address
	deviceConn.ld.TraceRouteList.PushBack()
	for e := newList.Front(); e != nil; e = e.Next() {
		fmt.Println(e.Value)
	}*/
	deviceConn.ld.TraceRouteList.PushBack(msg[1:])
	//TODO send request to the database to write the new IP list
}

/*
 * This method is called when a message's OP code is set to the KEYLOG_GET
 * constant. The new keylog file is then parsed in. A request is then sent
 * to the database to update with the new keylog data.
 */
func UpdateKeylog(deviceConn *deviceConnection, msg string) {
	deviceConn.ld.KeylogData.PushBack(msg[1 : len(msg)-1])
	fmt.Println(deviceConn.ld.KeylogData.Back().Value)
	if deviceConn.ld.UpdateKeylog() {
		fmt.Println("Keylog data has been successfully updated")
	} else {
		fmt.Println("Keylog data has NOT been successfully updated")
	}
	//TODO send request to database to add new keylog entry
}

/*
 * This method is always called immediately after a new connection is created.
 * The first thing a laptop should send whenever it connects is its ID (MAC
 * Address) and this is where it is read in. The connection object is then
 * hashed using the MAC Address.
 */
func GetDeviceID(conn net.Conn) { //(string, error) {
	buffer := make([]byte, 10240)
	//ld := new(LaptopDevice)
	bytesRead, err := conn.Read(buffer)
	if err != nil {
		fmt.Println("Error reading device ID", err)
	}
	//ld.ID = string(buffer[0:bytesRead])
	deviceConn := new(deviceConnection)
	deviceConn.ld.ID = string(buffer[0:bytesRead])
	deviceConn.conn = conn
	lh.mapDeviceQueue <- deviceConn
	var output string
	if deviceConn.ld.CheckIfStolen() {
		//TODO send device stolen OP code
		fmt.Println("CheckIfStolen request returned true")
		output = string(STOLEN) + "\n"
	} else {
		//TODO send device NOT stolen OP code
		fmt.Println("CheckIfStolen request returned false")
		output = string(NOT_STOLEN) + "\n"
	}
	outputBytes := []byte(output)
	bytesWritten, err := conn.Write(outputBytes)
	if err != nil {
		fmt.Println("Error sending stolen code.", err)
	}
	bytesWritten = bytesWritten
	//TODO have GetMessage be called in response to sending messages
	//go GetMessage(deviceConn)
}

/*
 * This method is where a laptop's open connection is hashed to its MAC Address
 * after the MAC Address (device ID) is read in in the GetDeviceID thread. This
 * method runs in its own thread because must wait on its channel to be filled
 * before running the hash, so most the time it is blocking the thread.
 */
func MapDeviceID() {
	for {
		dc := <-lh.mapDeviceQueue
		lh.connections[dc.ld.ID] = dc.conn
		fmt.Println(dc.ld.ID)
	}
}
