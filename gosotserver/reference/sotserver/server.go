package main

import (
	"bytes"
	"container/list"
	"fmt"
	"net"
)

func Log(v ...interface{}) {
	fmt.Println(v...)
}

func IOHandler(Incoming <-chan string, clientList *list.List) {
	for {
		Log("IOHandler: Waiting for input")
		input := <-Incoming
		Log("IOHandler: Handling ", input)
		for e := clientList.Front(); e != nil; e = e.Next() {
			client := e.Value.(Client)
			client.Incoming <- input
		}
	}
}

func ClientReader(client *Client) {
	buffer := make([]byte, 2048)

	for client.Read(buffer) {
		// Message for client to request disconnect
		// TODO: implement our own request types
		if bytes.Equal(buffer, []byte("/quit")) {
			client.Close()
			break
		}
		Log("ClientReader received ", client.id, "> ", string(buffer))
		//send := client.id + "> " + string(buffer)
		//client.Outgoing <-send
		// Flush buffer
		for i := 0; i < 2048; i++ {
			buffer[i] = 0x00
		}
	}

	client.Outgoing <- client.id + " has disconnected"
	Log("Client Reader stopped for ", client.id)
}

func ClientSender(client *Client) {
	for {
		select {
		case buffer := <-client.Incoming:
			Log("ClientSender sending ", string(2), " to ", client.id)
			count := 0
			for i := 0; i < len(buffer); i++ {
				if buffer[i] == 0x00 {
					break
				}
				count++
			}
			Log("Send size: ", count)
			client.Conn.Write([]byte(buffer)[0:count])
		case <-client.Quit:
			Log("Client ", client.id, " disconnecting")
			client.Conn.Close()
			break
		}
	}
}

func ClientHandler(conn net.Conn, ch chan string, clientList *list.List) {
	buffer := make([]byte, 1024)
	bytesRead, error := conn.Read(buffer)
	if error != nil {
		Log("Client connection error: ", error)
	}

	id := string(buffer[0:bytesRead])
	newClient := &Client{id, make(chan string), ch, conn, make(chan bool), clientList}

	go ClientSender(newClient)
	go ClientReader(newClient)
	clientList.PushBack(*newClient)
	ch <- string(2)
	//ch <- string(id + " has connected")
}
func UserInputReader(option int) {
	fmt.Print("Enter option: ")
	fmt.Scan(&option)
}

func main() {
	Log("Server started!")
	go startWebServer()
	Log("WebServer started!")

	var option int
	go UserInputReader(option)

	clientList := list.New()
	in := make(chan string)
	go IOHandler(in, clientList)

	service := ":10011"
	// Listen on TCP port 10011 on all interfaces
	l, err := net.Listen("tcp", service)
	if err != nil {
		Log(err)
	}
	defer l.Close()
	for {
		// Wait for a connection
		Log("Waiting for clients")
		conn, err := l.Accept()
		if err != nil {
			Log(err)
		} else {
			// Handle the connection in a new goroutine
			// The loop then returns to accepting, so that
			// multiple connections can be served concurrently
			go ClientHandler(conn, in, clientList)
		}
	}

}
