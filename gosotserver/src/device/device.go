/*
 * @author Nathan Plotts (nwp0002@auburn.edu)
 *
 * This file is where the structure for the parent device object is defined. Each
 * device whether, it is a GPS device or a laptop, will have the variables defined
 * in this struct.
 */

package device

//import "net"

type Device struct {
	ID string
}

/*type Client struct {
	Conn net.Conn
}*/

/*type DeviceBuffer struct {
	LaptopDevices map[string]LaptopDevice
	GPSDevices    map[string]GPSDevice
}*/
