package webserver

import (
	"CustomProtocol"
	"fmt"
)

// hub maintains the set of active connections and broadcasts messages to the
// connections.

type hub struct {
	// Registered connections
	connections map[string]*connection

	// Broadcast message
	broadcast chan []byte

	// Incoming messages from the connections
	inMessage chan Message

	// Out going message to connections
	outMessage chan []byte

	// Register requests from the connections
	register chan *connection

	// Unregister requests from connections
	unregister chan *connection
}

var h = hub{
	//inMessage: make(chan Message), // Send only channel
	broadcast:   make(chan []byte),
	outMessage:  make(chan []byte), // Receive only channel
	register:    make(chan *connection),
	unregister:  make(chan *connection),
	connections: make(map[string]*connection),
}

func updateClientMap() {

}

func (h *hub) run() {
	for {
		select {
		case c := <-h.register:
			for _, deviceId := range c.gpsDeviceList {
				h.connections[deviceId] = c
			}
			fmt.Println(h.connections)
		case c := <-h.unregister:
			for _, deviceId := range c.gpsDeviceList {
				delete(h.connections, deviceId)
			}
			close(c.send)
		case m := <-h.broadcast:
			// m is a []byte separated by 0x1B
			parsedPayload := CustomProtocol.ParsePayload(m)
			// NOTE: parsedPayload[0] == deviceId (Phone Number)
			fmt.Println(parsedPayload)
			msg := []byte{}
			// Latitude
			msg = append(msg, []byte(parsedPayload[1])...)
			msg = append(msg, 0x1B)
			// Longitude
			msg = append(msg, []byte(parsedPayload[2])...)
			h.connections[parsedPayload[0]].send <- msg
			/*for c := range h.connections {
				select {
				case c.send <- m:
				default:
					close(c.send)
					delete(h.connections, c)
				}
			}*/
		}
	}
}
