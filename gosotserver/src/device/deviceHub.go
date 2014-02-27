package device

import (
	"CustomRequest"
	"container/list"
	"fmt"
	"net"
	"strconv"
	"strings"
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

type deviceHub struct {
	connections map[string]net.Conn
	DeviceBuffer
	mapDeviceQueue chan *deviceConnection
}

type deviceConnection struct {
	ld   LaptopDevice
	conn net.Conn
}

var dh = deviceHub{
	mapDeviceQueue: make(chan *deviceConnection, 20),
	connections:    make(map[string]net.Conn),
}

var toServer chan *CustomRequest.Request
var fromServer chan *CustomRequest.Request
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
	dh.mapDeviceQueue <- deviceConn
	go GetMessage(deviceConn)
}

// Hash the device to the connection
func MapDeviceID() {
	for {
		dc := <-dh.mapDeviceQueue
		dh.connections[dc.ld.ID] = dc.conn
		fmt.Println(dc.ld.ID)
	}
}

func SmsConnection() {
	//connect
	listener, err := net.Listen(CONN_TYPE, CONN_PORT_SMS)
	if err != nil {
		fmt.Println("Error listening:", err.Error())
	} else {

	}

	fmt.Println("Connection created on " + CONN_TYPE + " " + CONN_PORT_SMS)

	//send & receive
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error connecting", err)
		}
		buffer := make([]byte, 512)
		fmt.Println("Connection established with SMS client")

		msg := ""
		for {
			bytesRead, _ := conn.Read(buffer)
			if bytesRead > 0 {
				if bytesRead > 10 {
					received := string(buffer[0:bytesRead])
					msg = googleMapLinkParser(received)
					fmt.Println("Received msg: ", msg)
					toServerT <- []byte(msg)
					//msg = msg + received
				} else {
					conn.Write([]byte("|"))
				}
			}
		}
	}

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

var toServerT chan []byte

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
