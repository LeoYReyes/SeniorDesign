package device

import (
	"fmt"
	"net"
)

const (
	CONN_TYPE = "tcp"
	CONN_PORT = ":10015"
)

type deviceHub struct {
	connections map[string]net.Conn
	DeviceBuffer
	mapDeviceQueue chan *mapDeviceQueueStruct
}

type mapDeviceQueueStruct struct {
	ld   LaptopDevice
	conn net.Conn
}

var dh = deviceHub{
	mapDeviceQueue: make(chan *mapDeviceQueueStruct, 100),
	connections:    make(map[string]net.Conn),
}

var deviceConn = new(mapDeviceQueueStruct)

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

// Get ID from device
func GetDeviceID(conn net.Conn) { //(string, error) {
	buffer := make([]byte, 1024)
	//ld := new(LaptopDevice)
	bytesRead, err := conn.Read(buffer)
	if err != nil {
		fmt.Println("Error reading device ID", err)
	}
	//ld.ID = string(buffer[0:bytesRead])
	deviceConn.ld.ID = string(buffer[0:bytesRead])
	deviceConn.conn = conn
	dh.mapDeviceQueue <- deviceConn
}

// Hash the device to the connection
func MapDeviceID() {
	fmt.Println("Before for loop")
	for {
		fmt.Println("Begin for loop")
		dc := <-dh.mapDeviceQueue
		fmt.Println("Before hashing")
		dh.connections[dc.ld.ID] = dc.conn
		fmt.Println("After hashing")
		fmt.Println(dc.ld.ID)
	}
}

// Read from open connections
func Read(conn net.Conn) bool {
	buffer := make([]byte, 1024)
	conn.Read(buffer)
	return false
}

// Test code for server.go
//go device.MapDeviceID()
//listener := device.Connect()
//device.Listen(listener)
