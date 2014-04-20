package webserver

import (
	"CustomProtocol"
	"fmt"
	"net/http"
	"strings"
)

func newDeviceHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		// Handle error
		fmt.Println("ParseForm error: ", err)
	}

	//fmt.Println("Form: ", r.PostForm)

	deviceType := strings.Trim(r.PostForm.Get("deviceType"), " ")
	deviceId := strings.Trim(r.PostForm.Get("deviceId"), " ")
	deviceName := strings.Trim(r.PostForm.Get("deviceName"), " ")

	// Check for any empty form inputs
	if deviceType == "" || deviceId == "" || deviceName == "" {
		// Send back error response, input field is empty

		return
	}

	buf := CustomProtocol.CreatePayload(deviceType, deviceId, deviceName)
	// Device owner, get user account info from session
	_, cookieErr := r.Cookie("userSession")
	if cookieErr != nil {
		fmt.Println("Cookie Error: ", cookieErr)
	} else {
		sesh, _ := store.Get(r, "userSession")
		/*buf = append(buf, []byte(sesh.Values["userId"].(string))...)
		buf = append(buf, 0x1B)*/
		buf = append(buf, CustomProtocol.CreatePayload(sesh.Values["userId"].(string))...)
		//fmt.Println("Cookie userid: ", sesh.Values["userId"])
		//fmt.Println("Cookie isLoggedIn: ", sesh.Values["isLoggedIn"])
	}
	//fmt.Println("buf: ", string(buf))
	// Create a response channel to receive response for the reqeust
	resCh := make(chan []byte)
	// Create request to register new device and send off request
	req := &CustomProtocol.Request{Id: CustomProtocol.AssignRequestId(), Destination: CustomProtocol.Database, Source: CustomProtocol.Web,
		OpCode: CustomProtocol.NewDevice, Payload: buf, Response: resCh}
	toServer <- req
	if deviceType == "gps" {
		//geogram setup
		geogramBuf := CustomProtocol.CreatePayload(deviceId, "1234") //todo hard-coded PIN for now
		response := make(chan []byte)
		geogramSetupReq := &CustomProtocol.Request{Id: CustomProtocol.AssignRequestId(), Destination: CustomProtocol.DeviceGPS, Source: CustomProtocol.Web,
			OpCode: CustomProtocol.GeogramSetup, Payload: geogramBuf, Response: response}
		toServer <- geogramSetupReq

		//geogram sleep
		response2 := make(chan []byte)
		geogramBuf2 := CustomProtocol.CreatePayload(deviceId, "1234") //todo hard-coded PIN for now
		geogramSetupReq2 := &CustomProtocol.Request{Id:               CustomProtocol.AssignRequestId(), Destination: CustomProtocol.DeviceGPS, Source: CustomProtocol.Web,
			OpCode: CustomProtocol.SleepGeogram, Payload: geogramBuf2, Response: response2}
		toServer <- geogramSetupReq2

		//fmt.Println(deviceId + " Geogram setup complete")
	}
	successful, res := CustomProtocol.GetResponse(resCh, 10)
	if successful {
		// response is true if successfully registered, false if there is an error
		//fmt.Println("Response: ", res[0])
		//TODO: notify client of the response
		if res[0] == 0 {
			fmt.Println("New Device Registration: Failed")
			w.Write([]byte("fail"))
		} else {
			fmt.Println("New Device Registration: Success")
			w.Write([]byte("success"))
		}

	} else {
		fmt.Println("No response on device registration")
		http.Error(w, "Timeout on device registration", 666)
		return
	}

}

// Toggles the device's stolen status
func toggleDeviceHandler(w http.ResponseWriter, r *http.Request) {

	err := r.ParseForm()
	if err != nil {
		// Handle error
		fmt.Println("ParseForm error: ", err)
	}

	//fmt.Println("Form: ", r.PostForm)
	//TODO: user input for geogram PIN
	deviceType := strings.Trim(r.PostForm.Get("deviceType"), " ")
	deviceId := strings.Trim(r.PostForm.Get("deviceId"), " ")
	deviceCommand := strings.Trim(r.PostForm.Get("deviceCommand"), " ")

	// Check for any empty form inputs
	if deviceType == "" || deviceId == "" || deviceCommand == "" {
		// Send back error response, input field is empty

		return
	}
	resCh := make(chan []byte)
	buf := CustomProtocol.CreatePayload(deviceId)
	//fmt.Println("Device type: ", deviceType)

	switch deviceType {
	case "gps":
		//todo Default PIN-NUMBER for Geogram One

		buf = append(buf, CustomProtocol.CreatePayload("1234")...)
		if deviceCommand == "1" {
			reqToDB := &CustomProtocol.Request{Id: CustomProtocol.AssignRequestId(), Destination: CustomProtocol.Database, Source: CustomProtocol.Web,
				OpCode: CustomProtocol.ActivateGPS, Payload: buf, Response: resCh}
			toServer <- reqToDB
			// Default interval 140 seconds
			buf = append(buf, CustomProtocol.CreatePayload("140")...)
			reqToDevice := &CustomProtocol.Request{Id: CustomProtocol.AssignRequestId(), Destination: CustomProtocol.DeviceGPS, Source: CustomProtocol.Web,
				OpCode: CustomProtocol.ActivateIntervalGps, Payload: buf, Response: nil}
			toServer <- reqToDevice
			//sleep command
			reqToDevice2 := &CustomProtocol.Request{Id: CustomProtocol.AssignRequestId(), Destination: CustomProtocol.DeviceGPS, Source: CustomProtocol.Web,
				OpCode: CustomProtocol.SleepGeogram, Payload: buf, Response: nil}
			toServer <- reqToDevice2
		} else {
			// send to database to flag not stolen
			reqToDB := &CustomProtocol.Request{Id: CustomProtocol.AssignRequestId(), Destination: CustomProtocol.Database, Source: CustomProtocol.Web,
				OpCode: CustomProtocol.FlagNotStolen, Payload: buf, Response: resCh}
			toServer <- reqToDB
			// Deactivate command to device
			// end tracking
			buf = append(buf, CustomProtocol.CreatePayload("0")...)
			reqToDevice := &CustomProtocol.Request{Id: CustomProtocol.AssignRequestId(), Destination: CustomProtocol.DeviceGPS, Source: CustomProtocol.Web,
				OpCode: CustomProtocol.ActivateIntervalGps, Payload: buf, Response: nil}
			toServer <- reqToDevice
		}
	case "laptop":
		if deviceCommand == "1" {
			req := &CustomProtocol.Request{Id: CustomProtocol.AssignRequestId(), Destination: CustomProtocol.Database, Source: CustomProtocol.Web,
				OpCode: CustomProtocol.FlagStolen, Payload: buf, Response: resCh}
			toServer <- req
		} else {
			// send to database to flag not stolen
			reqToDB := &CustomProtocol.Request{Id: CustomProtocol.AssignRequestId(), Destination: CustomProtocol.Database, Source: CustomProtocol.Web,
				OpCode: CustomProtocol.FlagNotStolen, Payload: buf, Response: resCh}
			toServer <- reqToDB

			reqToDeviceHub := &CustomProtocol.Request{Id: CustomProtocol.AssignRequestId(), Destination: CustomProtocol.DeviceLaptop, Source: CustomProtocol.Web,
				OpCode: CustomProtocol.FlagNotStolen, Payload: buf, Response: resCh}
			toServer <- reqToDeviceHub
		}
	default:
	}

}

func pingDeviceHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		// Handle error
		fmt.Println("ParseForm error: ", err)
	}

	deviceType := strings.Trim(r.PostForm.Get("deviceType"), " ")
	deviceId := strings.Trim(r.PostForm.Get("deviceId"), " ")
	deviceCommand := strings.Trim(r.PostForm.Get("deviceCommand"), " ")

	// Check for any empty form inputs
	if deviceType == "" || deviceId == "" || deviceCommand == "" {
		// Send back error response, input field is empty

		return
	}

	resCh := make(chan []byte)
	buf := []byte{}
	buf = CustomProtocol.CreatePayload(deviceId, deviceCommand)

	req := &CustomProtocol.Request{Id: CustomProtocol.AssignRequestId(), Destination: CustomProtocol.DeviceGPS, Source: CustomProtocol.Web,
		OpCode: CustomProtocol.FreestyleMsg, Payload: buf, Response: resCh}
	toServer <- req
}

func deviceInfoHandler(w http.ResponseWriter, r *http.Request) {
	buf := []byte{}
	// Device owner, get user account info from session
	_, cookieErr := r.Cookie("userSession")
	if cookieErr != nil {
		fmt.Println("Cookie Error: ", cookieErr)
	} else {
		sesh, _ := store.Get(r, "userSession")
		/*buf = append(buf, []byte(sesh.Values["userId"].(string))...)
		buf = append(buf, 0x1B)*/
		buf = CustomProtocol.CreatePayload(sesh.Values["userId"].(string))
		//fmt.Println("Cookie userid: ", sesh.Values["userId"])
		//fmt.Println("Cookie isLoggedIn: ", sesh.Values["isLoggedIn"])
	}
	// Create a response channel to receive response for the reqeust
	resCh := make(chan []byte)
	// Create request to register new device and send off request
	req := &CustomProtocol.Request{Id: CustomProtocol.AssignRequestId(), Destination: CustomProtocol.Database, Source: CustomProtocol.Web,
		OpCode: CustomProtocol.GetDeviceList, Payload: buf, Response: resCh}
	toServer <- req
	successful, res := CustomProtocol.GetResponse(resCh, 10)
	if successful {
		str := string(res)
		finalRes := []byte{}
		resLaptop := str[:strings.Index(str, string(0x1B))]
		resGPS := str[strings.Index(str, string(0x1B))+1:]
		//fmt.Println(resLaptop)
		//fmt.Println(resGPS)
		if len(resLaptop) > 5 {
			if len(resGPS) > 5 {
				resLaptop = str[:strings.Index(str, string(0x1B))-1]
				resGPS = str[strings.Index(str, string(0x1B))+2:]
			}
			finalRes = append(finalRes, []byte(resLaptop)...)
			w.Header().Set("Content-Type", "application/json")
		}
		if len(resGPS) > 5 {
			if len(resLaptop) > 5 {
				finalRes = append(finalRes, 0x2C)
			}
			finalRes = append(finalRes, []byte(resGPS)...)
			w.Header().Set("Content-Type", "application/json")
		}
		w.Write(finalRes)
	} else {
		return
	}

}
