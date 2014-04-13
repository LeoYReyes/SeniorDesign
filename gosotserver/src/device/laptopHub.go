/*f
 * @author Nathan Plotts (nwp0002@auburn.edu)
 * This file is where the the laptopHub struct is defined and
 * also where general laptop connection handling functions will
 * be stored.
 */

package device

import (
	"CustomProtocol"
	//"container/list"
	"fmt"
	"net"
	//"strconv"
	"strings"
	//"time"
)

type laptopHub struct {
	connections map[string]net.Conn
	//	DeviceBuffer
	//mapDeviceQueue chan *deviceConnection
}

type deviceConnection struct {
	ld   LaptopDevice
	conn net.Conn
}

var lh = laptopHub{
	//mapDeviceQueue: make(chan *deviceConnection, 20),
	connections: make(map[string]net.Conn),
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
func GetMessage(deviceConn deviceConnection) {
	defer func() {
		CloseConn(deviceConn)
	}()
	buffer := make([]byte, 10240)
	badMsg := 0
	for {
		fmt.Println("Waiting for message from client...")
		_, err := deviceConn.conn.Read(buffer)
		if err != nil {
			fmt.Println("Error reading", err)
			break
		}
		msg := string(buffer)
		index := strings.Index(msg, "\n")
		if index != -1 {
			msg = msg[1:index]
		} else {
			fmt.Println("LaptopDevice.GetDeviceID: deviceId Parse Error")
		}
		//opCode, err := strconv.Atoi(msg[0:1])
		opCode := buffer[0]
		//fmt.Println("Message byte format: ", buffer[1:bytesRead])
		fmt.Println("Message received with OP Code: ", opCode)
		if err != nil {
			fmt.Println("Invalid OP code", err)
		} else {
			success := true
			switch opCode {
			case CustomProtocol.UpdateUserIPTraceData:
				UpdateTraceroute(deviceConn, msg)
			case CustomProtocol.UpdateUserKeylogData:
				UpdateKeylog(deviceConn, msg)
			case CustomProtocol.NoOp:
			default:
				success = false
			}

			if success {
				badMsg = 0
			} else {
				badMsg += 1
				if badMsg > 4 {
					break
				}
			}
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
func UpdateTraceroute(deviceConn deviceConnection, msg string) {
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
	fmt.Println(msg)
	ipAddr := deviceConn.conn.RemoteAddr().String()
	msgBytes := []byte{}
	msgBytes = append(msgBytes, []byte(ipAddr)...)
	msgBytes = append(msgBytes, 0x7E)
	msgBytes = append(msgBytes, []byte(msg)...)
	msg = string(msgBytes)
	deviceConn.ld.TraceRouteList = append(deviceConn.ld.TraceRouteList, msg)
	if deviceConn.ld.UpdateTraceroute() {
		fmt.Println("Traceroute data has been successfully updated")
	} else {
		fmt.Println("Traceroute data has NOT been successfully updated")
	}
	//TODO send request to the database to write the new IP list
}

/*
 * This method is called when a message's OP code is set to the KEYLOG_GET
 * constant. The new keylog file is then parsed in. A request is then sent
 * to the database to update with the new keylog data.
 */
func UpdateKeylog(deviceConn deviceConnection, msg string) {
	deviceConn.ld.KeylogData = append(deviceConn.ld.KeylogData, msg) //[1:len(msg)-1])
	fmt.Println(deviceConn.ld.KeylogData[len(deviceConn.ld.KeylogData)-1])
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
	_, err := conn.Read(buffer)
	if err != nil {
		fmt.Println("Error reading device ID", err)
	}
	//ld.ID = string(buffer[0:bytesRead])
	deviceConn := new(deviceConnection)
	deviceId := string(buffer)
	index := strings.Index(deviceId, "\n")
	if index != -1 {
		deviceConn.ld.ID = deviceId[:index]
	} else {
		fmt.Println("LaptopDevice.GetDeviceID: deviceId Parse Error")
		deviceConn.ld.ID = ""
	}

	deviceConn.conn = conn
	MapDeviceID(deviceConn)
	var sentStolen bool
	if deviceConn.ld.CheckIfStolen() {
		fmt.Println("CheckIfStolen request returned true")
		sentStolen = SendMsg(deviceConn.ld.ID, CustomProtocol.FlagStolen, "")
		// SEND requests to laptop upon connection, if stolen
		requestTraceRoute := SendMsg(deviceConn.ld.ID, CustomProtocol.UpdateUserIPTraceData, "")
		requestKeyLog := SendMsg(deviceConn.ld.ID, CustomProtocol.UpdateUserKeylogData, "")
		if !sentStolen {
			fmt.Println("Error sending stolen code.")
		}
		if !requestTraceRoute {
			fmt.Println("laptopHub.GetDeviceID: Error sending request traceroute")
		}
		if !requestKeyLog {
			fmt.Println("laptopHub.GetDeviceID: Eror sending request keylog")
		}
		go GetMessage(*deviceConn)
	} else { //if CheckIfStolen returns false
		fmt.Println("CheckIfStolen request returned false")
		ipAddr := conn.RemoteAddr()
		fmt.Println(ipAddr)
		sentStolen = SendMsg(deviceConn.ld.ID, CustomProtocol.FlagNotStolen, "")
		if !sentStolen {
			fmt.Println("Error sending stolen code.")
		}
		CloseConn(*deviceConn)
		/*err := conn.Close()
		if err != nil {
			fmt.Println("Error closing laptop connection.", err)
		}*/
		fmt.Println("Connection sucks-s-foli closed")
		//todo close connection and laptop goes into sleep mode
	}
	//TODO have GetMessage be called in response to sending messages
	//go GetMessage(deviceConn)
}

/*
 * This method sends a message to a laptop if a connection to it is found.
 * It uses the laptop's ID (MAC address) to search for the connection in the map
 * of connections, and sends a message in the format <opcode><payload>
 */
func SendMsg(id string, opcode byte, payload string) bool {
	conn := lh.connections[id]
	if conn == nil {
		fmt.Println("SendMsg: Connection not found for ID " + id)
		return false
	}
	var op [1]byte
	op[0] = opcode
	msg := append(op[0:1], []byte(payload)...)
	_, err := conn.Write(msg)
	if err != nil {
		fmt.Println("SendMsg: Error sending message to device with ID " + id)
		return false
	}
	fmt.Printf("Message %d"+payload+" sent to device with ID "+id+"\n", opcode)
	return true
}

/*
 * This method will process laptop related requests and return true if the
 * message is sent to the laptop
 */
func ProcessLapReq(req *CustomProtocol.Request) {
	id := string(req.Payload)
	sent := SendMsg(id, req.OpCode, "")
	if sent {
		req.Response <- []byte{1}
	} else {
		req.Response <- []byte{0}
	}
}

/*
 * This method is where a laptop's open connection is hashed to its MAC Address
 * after the MAC Address (device ID) is read in in the GetDeviceID thread. This
 * method runs in its own thread because must wait on its channel to be filled
 * before running the hash, so most the time it is blocking the thread.
 */
/*
func MapDeviceID() {
	for {
		dc := <-lh.mapDeviceQueue
		lh.connections[dc.ld.ID] = dc.conn
		fmt.Println(dc.ld.ID)
	}
}
*/

/*
 * Adds a connection to the connections map keyed by the ID (MAC address)
 * of the device connecting
 */
func MapDeviceID(dc *deviceConnection) {
	if lh.connections[dc.ld.ID] != nil { //todo newly added, make sure it is safe
		lh.connections[dc.ld.ID].Close()
		lh.connections[dc.ld.ID] = nil
	}
	lh.connections[dc.ld.ID] = dc.conn
	fmt.Println(dc.ld.ID)
}

/*
 * Closes a connection and removes it from the map
 */
func CloseConn(dc deviceConnection) {
	dc.conn.Close()
	lh.connections[dc.ld.ID] = nil
	fmt.Println(dc.ld.ID + ": connection closed and removed")
}
