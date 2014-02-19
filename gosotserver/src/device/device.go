package device

import "net"

type Device struct {
	ID string
}

type Client struct {
	Conn net.Conn
}

type DeviceBuffer struct {
	LaptopDevices map[string]LaptopDevice
	GPSDevices    map[string]GPSDevice
}
