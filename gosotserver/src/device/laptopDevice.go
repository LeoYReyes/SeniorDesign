/*
 * @author Nathan Plotts (nwp0002@auburn.edu)
 * This file is where the the LaptopDevice struct is defined and
 * also where LaptopDevice specific functions are stored.
 */

package device

import (
	"CustomProtocol"
	//"container/list"
	"fmt"
)

/*
 * This struct defines the structure of a connected laptop device. It contains a list
 * of strings that contain IP traceroutes and a list of strings that contain keylog
 * data.
 */
type LaptopDevice struct {
	TraceRouteList []string
	// List of Key Logs
	KeylogData []string
	Device
}

/*type LaptopClient struct {
	Client
}*/

/*
 * This function creates and sends a request to the database to check if the connected
 * laptop is stolen. If it is stolen it sends a message back to the laptop alerting it
 * that it has been marked as stolen and needs to start tracking. If the request says
 * that the laptop is not stolen it sends a message to the laptop letting it know that
 * it is not stolen.
 */
func (ld *LaptopDevice) CheckIfStolen() bool {
	//TODO send database request here
	id := []byte(ld.Device.ID)
	id = append(id, 0x1B)
	fmt.Println(id)
	response := make(chan []byte)
	req := &CustomProtocol.Request{Id: CustomProtocol.AssignRequestId(), Destination: CustomProtocol.Database,
		Source: CustomProtocol.DeviceLaptop, OpCode: CustomProtocol.CheckDeviceStolen, Payload: id,
		Response: response}
	toServer <- req
	fmt.Println("CheckIfStolen request created and sent")
	isStolen := <-response
	if isStolen[0] == 1 {
		return true
	} else {
		return false
	}
}

/*
 * This function creates and sends a request to the database to write new keylog data
 * to the device's logs. The database function would then respond with boolean letting
 * the function know if the write succeeded.
 */
func (ld *LaptopDevice) UpdateKeylog() bool {
	id := []byte(ld.Device.ID)
	keylog := ld.KeylogData[len(ld.KeylogData)-1]
	payload := append(id, 0x1B)
	payload = append(payload, keylog...)
	payload = append(payload, 0x1B)
	response := make(chan []byte)
	req := &CustomProtocol.Request{Id: CustomProtocol.AssignRequestId(), Destination: CustomProtocol.Database,
		Source: CustomProtocol.DeviceLaptop, OpCode: CustomProtocol.UpdateUserKeylogData, Payload: payload,
		Response: response}
	toServer <- req
	fmt.Println("UpdateKeylogData request created and sent")
	KeylogUpdated := <-response
	if KeylogUpdated[0] == 1 {
		return true
	} else {
		return false
	}
}

/*
 * This function creates and sends a request to the database to write a new IP traceroute
 * string to the device's logs. The database then response with a boolean to let the
 * function know if the write succeeded.
 */
func (ld *LaptopDevice) UpdateTraceroute() bool {
	id := []byte(ld.Device.ID)
	traceroute := ld.TraceRouteList[len(ld.TraceRouteList)-1]
	tracerouteBytes := []byte(traceroute)
	payload := append(id, 0x1B)
	payload = append(payload, tracerouteBytes...)
	payload = append(payload, 0x1B)
	response := make(chan []byte)
	req := &CustomProtocol.Request{Id: CustomProtocol.AssignRequestId(), Destination: CustomProtocol.Database,
		Source: CustomProtocol.DeviceLaptop, OpCode: CustomProtocol.UpdateUserIPTraceData, Payload: payload,
		Response: response}
	toServer <- req
	fmt.Println("UpdateTracerouteData request created and sent")
	TracerouteUpdated := <-response
	if TracerouteUpdated[0] == 1 {
		return true
	} else {
		return false
	}
	//TODO possibly only deserialize the IP list when displaying on the webpage
}
