/*
 * @author Nathan Plotts (nwp0002@auburn.edu)
 * This file is where the the LaptopDevice struct is defined and
 * also where LaptopDevice specific functions will be stored.
 */

package device

import "container/list"
import "CustomProtocol"

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
	isStolen := <-response
	if isStolen[0] == 1 {
		return true
	} else {
		return false
	}
}
