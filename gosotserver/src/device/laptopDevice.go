/*
 * @author Nathan Plotts (nwp0002@auburn.edu)
 * This file is where the the LaptopDevice struct is defined and
 * also where LaptopDevice specific functions will be stored.
 */

package device

import "container/list"

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
	ld = ld
	return true
}
