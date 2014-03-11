package webserver

import (
	"fmt"
	"net/http"
)

func newDeviceHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		// Handle error
		fmt.Println("ParseForm error: ", err)
	}

	fmt.Println(r.PostForm)
	buf := []byte{}
	// Device Id (phone number for Geogram, MAC Address for laptops)
	buf = append(buf, []byte(r.PostForm.Get("deviceId"))...)
	buf = append(buf, 0x1B)
	// Device type
	buf = append(buf, []byte(r.PostForm.Get("deviceType"))...)
	buf = append(buf, 0x1B)
	// Device name
	buf = append(buf, []byte(r.PostForm.Get("deviceName"))...)
	buf = append(buf, 0x1B)
	// Device owner, get user account info from session
	session, _ := store.Get(r, "userSession")
	buf = append(buf, []byte(session.Values["userId"])...) // NEEDS FIX
	buf = append(buf, 0x1B)

	if err != nil {
		// Handle error
		fmt.Println("decoder error: ", err)
	}
}
