package device

import "container/list"

type LaptopDevice struct {
	TraceRouteList list.List
	Device
}

type LaptopClient struct {
	Client
}
