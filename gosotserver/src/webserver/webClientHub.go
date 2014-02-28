package webserver

// hub maintains the set of active connections and broadcasts messages to the
// connections.

type hub struct {
	// Registered connections
	connections map[*connection]bool

	// Broadcast message
	broadcast chan []byte

	// Incoming messages from the connections
	inMessage chan Message

	// Out going message to connections
	outMessage chan Message

	// Register requests from the connections
	register chan *connection

	// Unregister requests from connections
	unregister chan *connection
}

var h = hub{
	//inMessage: make(chan Message), // Send only channel
	broadcast:   make(chan []byte),
	outMessage:  make(chan Message), // Receive only channel
	register:    make(chan *connection),
	unregister:  make(chan *connection),
	connections: make(map[*connection]bool),
}

func updateClientMap() {

}

func (h *hub) run() {
	for {
		select {
		case c := <-h.register:
			h.connections[c] = true
		case c := <-h.unregister:
			delete(h.connections, c)
			close(c.send)
		case m := <-h.broadcast:
			for c := range h.connections {
				select {
				case c.send <- m:
				default:
					close(c.send)
					delete(h.connections, c)
				}
			}
		}
	}
}
