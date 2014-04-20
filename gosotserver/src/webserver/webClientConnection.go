package webserver

import (
	"CustomProtocol"
	//"encoding/gob"
	"device"
	"encoding/json"
	"fmt"
	"net/http"
	"sessions"
	"time"
	"websocket"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

// Not being used right now, might not need at all
type userSession struct {
	userId      string
	inSession   bool
	isAdmin     bool
	currentPage string
}

// connection is an middleman between the websocket connection and the hub
type connection struct {
	// The websocket connection
	ws *websocket.Conn

	// GPS Device IDs associated with connection
	gpsDeviceList []string

	// Buffered channel of outbound messages
	send chan []byte
}

type Message struct {
	conn    *connection
	message []byte
}

// readPump pumps messages from the websocket connection to the hub
func (c *connection) readPump() {
	defer func() {
		h.unregister <- c
		c.ws.Close()
	}()
	c.ws.SetReadLimit(maxMessageSize)
	c.ws.SetReadDeadline(time.Now().Add(pongWait))
	c.ws.SetPongHandler(func(string) error { c.ws.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		_, message, err := c.ws.ReadMessage()
		if err != nil {
			break
		}
		//fmt.Println(string(message))
		// Creates request from message and send off
		// Currently only sending msg to GPS device
		response := make(chan []byte)
		req := &CustomProtocol.Request{Id: CustomProtocol.AssignRequestId(), Destination: CustomProtocol.DeviceGPS,
			Source: CustomProtocol.Web, OpCode: CustomProtocol.ActivateGPS, Payload: message, Response: response}
		toServer <- req
		//fmt.Println("Response received: ", <-response)
	}
}

// write writes a message with the given message type and payload
func (c *connection) write(mt int, payload []byte) error {
	c.ws.SetWriteDeadline(time.Now().Add(writeWait))
	return c.ws.WriteMessage(mt, payload)
}

// writePump pumps messages from the hub to the websocket connection
func (c *connection) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.ws.Close()
	}()
	for {
		select {
		case message, ok := <-c.send:
			if !ok {
				c.write(websocket.CloseMessage, []byte{})
				return
			}
			if err := c.write(websocket.TextMessage, message); err != nil {
				return
			}
		case <-ticker.C:
			if err := c.write(websocket.PingMessage, []byte{}); err != nil {
				return
			}
		}
	}
}

// serverWs handles webocket requests from the peer
func serveWs(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", 405)
		return
	}
	if r.Header.Get("Origin") != "http://"+r.Host {
		http.Error(w, "Origin not allowed", 403)
		return
	}
	ws, err := websocket.Upgrade(w, r, nil, 1024, 1024)
	if _, ok := err.(websocket.HandshakeError); ok {
		http.Error(w, "Not a websocket handshake", 400)
		return
	} else if err != nil {
		fmt.Println(err)
		return
	}
	//TODO: see user info from r (*http.Request) get deviceID list from session
	_, cookieErr := r.Cookie("userSession")
	sesh, _ := store.Get(r, "userSession")
	if cookieErr != nil {
		fmt.Println("Cookie Error: ", cookieErr)
	} else {

	}
	userName := string(sesh.Values["userId"].(string))

	resCh := make(chan []byte)
	/*buf := []byte{}
	buf = append(buf, []byte("gps")...)
	buf = append(buf, 0x1B)
	buf = append(buf, []byte(userName)...)
	buf = append(buf, 0x1B)*/
	buf := CustomProtocol.CreatePayload("gps", userName)
	req := &CustomProtocol.Request{Id: CustomProtocol.AssignRequestId(), Destination: CustomProtocol.Database, Source: CustomProtocol.Web,
		OpCode: CustomProtocol.GetDevice, Payload: buf, Response: resCh}

	toServer <- req
	res := <-resCh
	// res comes in as JSON data
	var deviceList []device.Device
	err = json.Unmarshal(res, &deviceList)
	if err != nil {
		fmt.Println("deviceList unmarshal error")
	}
	//list := CustomProtocol.ParsePayload(res)
	fmt.Println("serveWS list: ", deviceList)
	list := []string{}
	for _, device := range deviceList {
		list = append(list, device.ID)
	}

	c := &connection{send: make(chan []byte, 256), ws: ws, gpsDeviceList: list}
	fmt.Println("serveWS conn: ", c)
	h.register <- c
	go c.writePump()
	c.readPump()
}

//TODO: add key rotation
var store = sessions.NewCookieStore([]byte("its-the-most-wonderful-time"))

func serveSession(w http.ResponseWriter, r *http.Request) {
	// Get a session. We're ignoring the error resulted from decoding an
	// existing session: Get() always returns a session, even if empty.
	session, _ := store.Get(r, "userSession")
	if session.IsNew {
		//fmt.Println("New session created")
		session.Values["userId"] = r.PostForm.Get("loginName")
		// Check to see if user is an admin
		//session.Values["isAdmin"] = true
	} else {
		//fmt.Println("Existing session loaded")
	}
	session.Values["isLoggedIn"] = "true"
	//TODO: Request database for device IDs associated with account
	//		create a Request to be sent to database
	//req := CustomRequest.Request{0, 1, 2, CustomRequest.GetDeviceList, "test"}
	//toServer <- &req
	session.Options = &sessions.Options{
		// MaxAge=0 means no 'Max-Age' attribute specified.
		// MaxAge<0 means delete cookie now, equivalently 'Max-Age: 0'.
		// MaxAge>0 means Max-Age attribute present and given in seconds.
		Path:     "/",
		MaxAge:   86400,
		HttpOnly: true,
	}
	//fmt.Println("Session store: ", store.Codecs)
	// Save it.
	err := session.Save(r, w)
	if err != nil {
		fmt.Println("Session save error")
	} else {
		http.Redirect(w, r, "/UserMapNEW", http.StatusFound)
	}
}
