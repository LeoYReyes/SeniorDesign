/*
 * @author Nathan Plotts (nwp0002@auburn.edu)
 * This file is where the the LaptopDevice struct is defined and
 * also where LaptopDevice specific functions will be stored.
 */

package device

import (
	"CustomProtocol"
	"container/list"
	"fmt"
)

type LaptopDevice struct {
	TraceRouteList list.List
	KeylogData     list.List
	Device
}

type LaptopClient struct {
	Client
}

func (ld *LaptopDevice) CheckIfStolen() bool {
	//TODO send database request here
	id := []byte(ld.Device.ID)
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

func (ld *LaptopDevice) UpdateKeylog() bool {
	id := []byte(ld.Device.ID)
	keylog := ld.KeylogData.Back().Value.(string)
	payload := append(id, keylog...)
	response := make(chan []byte)
	req := &CustomProtocol.Request{Id: CustomProtocol.AssignRequestId(), Destination: CustomProtocol.Database,
		Source: CustomProtocol.DeviceLaptop, OpCode: CustomProtocol.CheckDeviceStolen, Payload: payload,
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

func (ld *LaptopDevice) UpdateTraceroute() bool {
	id := []byte(ld.Device.ID)
	traceroute := ld.TraceRouteList.Back().Value.(string)
	tracerouteBytes := []byte(traceroute)
	payload := append(id, tracerouteBytes...)
	response := make(chan []byte)
	req := &CustomProtocol.Request{Id: CustomProtocol.AssignRequestId(), Destination: CustomProtocol.Database,
		Source: CustomProtocol.DeviceLaptop, OpCode: CustomProtocol.CheckDeviceStolen, Payload: payload,
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
