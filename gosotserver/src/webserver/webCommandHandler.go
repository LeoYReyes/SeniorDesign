package webserver

import (
	"CustomProtocol"
	"fmt"
	"net/http"
)

func newDeviceHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		// Handle error
		fmt.Println("ParseForm error: ", err)
	}

	fmt.Println("Form: ", r.PostForm)
	buf := []byte{}
	// Device type
	buf = append(buf, []byte(r.PostForm.Get("deviceType"))...)
	buf = append(buf, 0x1B)
	// Device Id (phone number for Geogram, MAC Address for laptops)
	buf = append(buf, []byte(r.PostForm.Get("deviceId"))...)
	buf = append(buf, 0x1B)
	// Device name
	buf = append(buf, []byte(r.PostForm.Get("deviceName"))...)
	buf = append(buf, 0x1B)
	// Device owner, get user account info from session
	_, cookieErr := r.Cookie("userSession")
	if cookieErr != nil {
		fmt.Println("Cookie Error: ", cookieErr)
	} else {
		sesh, _ := store.Get(r, "userSession")
		buf = append(buf, []byte(sesh.Values["userId"].(string))...)
		buf = append(buf, 0x1B)
		//fmt.Println("Cookie userid: ", sesh.Values["userId"])
		//fmt.Println("Cookie isLoggedIn: ", sesh.Values["isLoggedIn"])
	}
	fmt.Println("buf: ", string(buf))
	// Create a response channel to receive response for the reqeust
	resCh := make(chan []byte)
	// Create request to register new device and send off request
	req := &CustomProtocol.Request{Id: CustomProtocol.AssignRequestId(), Destination: CustomProtocol.Database, Source: CustomProtocol.Web,
		OpCode: CustomProtocol.NewDevice, Payload: buf, Response: resCh}
	toServer <- req
	res := <-resCh
	// response is true if successfully registered, false if there is an error
	fmt.Println("Response: ", res[0])
	//TODO: notify client of the response
	if res[0] == 0 {
		fmt.Println("New Device Registration: Failed")
		fmt.Fprintf(w, "failed")
	} else {
		fmt.Fprintf(w, "success")
		fmt.Println("New Device Registration: Success")
	}

}

func toggleDeviceHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		// Handle error
		fmt.Println("ParseForm error: ", err)
	}
}
