/*
 * @author Nathan Plotts (nwp0002@auburn.edu)
 * This file is where the structure for a GPSDevice is defined. Any functions
 * specific to the GPSDevice will also be defined here.
 */

package device

import (
	"CustomProtocol"
	//"fmt"
	"time"
)

const (
	HYPERLINK_1    = "http://maps.google.com/maps?q="
	HYPERLINK_2    = "+("
	HYPERLINK_3    = ")&z=19"
	MOTION_ALERT   = "motion alert!"
	GEOFENCE_ALERT = "left geofence!"
)

type GPSDevice struct {
	Coordinates []string
	Device
}

/*
 * Sends update request to webserver and database for gps coords. param should
 * ID (phone number), lat, long delimited by escape chars
 */
func UpdateMapCoords(payload string) {
	webResponse := false
	databaseResponse := false
	webResponseChan := make(chan []byte)
	reqWeb := &CustomProtocol.Request{Id: CustomProtocol.AssignRequestId(), Destination: CustomProtocol.Web, Source: CustomProtocol.DeviceGPS,
		OpCode: CustomProtocol.UpdateWebMap, Payload: []byte(payload), Response: webResponseChan}
	toServer <- reqWeb
	//fmt.Println("gpsDevice: Req sent to server")

	dataResponseChan := make(chan []byte)
	reqData := &CustomProtocol.Request{Id: CustomProtocol.AssignRequestId(), Destination: CustomProtocol.Database, Source: CustomProtocol.DeviceGPS,
		OpCode: CustomProtocol.UpdateDeviceGPS, Payload: []byte(payload), Response: dataResponseChan}
	toServer <- reqData
	//fmt.Println("gpsDevice: Req sent to database")
	// todo enable this if we need to wait for a response and add handling here
	for !webResponse /* || !databaseResponse */ { /*
			if !webResponse {
				select {
				case webRet := <-reqWeb.Response:
					webResponse = true
					fmt.Println("gpsDevice: web server request returned ", webRet)
				default:
					time.Sleep(10000 * time.Millisecond)
				}
			}*/

		if !databaseResponse {
			select {
			case dataRet := <-reqWeb.Response:
				databaseResponse = true
				//fmt.Println("gpsDevice: database request returned ", dataRet)
			default:
				time.Sleep(10000 * time.Millisecond)
			}
		}
	}
}

/*type GPSClient struct {
	Client
}*/
