/* This is a number of tests for deviceHub. Right now, it currently tests the
 * ability to process gps requests. It requires a little bit of manual work
 * work to confirm some of the outputs. To run the tests, make sure a phone
 * is running the SMS Gateway app and ready to connect to when the tests begin.
 * The phoneNumber variable below determines what phone number the processed
 * requests will be sent to. The unit tests themselves only check that a
 * response is received from each processed request (responses are sent through
 * and channel that is passed in the request) and it is the expected value.
 */

package device

import (
	"CustomProtocol"
	"testing"
	"time"
)

// use a phone number to manually test if the output was correct
var phoneNumber = "Put Phone Number Here"
var pin = "1234"
var toDeviceCh = make(chan *CustomProtocol.Request)
var fromDeviceCh = make(chan *CustomProtocol.Request)

// clears to server channel so it doesn't remain full
func hack() {
	//var m *CustomProtocol.Request
	for {
		m := <-toServer
		m = m
	}
}

func TestFreestyleMsg(t *testing.T) {
	go StartDeviceServer(fromDeviceCh, toDeviceCh)
	go hack()

	time.Sleep(10000 * time.Millisecond)

	resCh := make(chan []byte)

	buf := []byte{}
	buf = append(buf, []byte(phoneNumber)...)
	buf = append(buf, 0x1B)
	buf = append(buf, []byte("Begin Tests :D")...)
	buf = append(buf, 0x1B)

	req := &CustomProtocol.Request{Id: CustomProtocol.AssignRequestId(), Destination: CustomProtocol.DeviceGPS, Source: CustomProtocol.Web,
		OpCode: CustomProtocol.FreestyleMsg, Payload: buf, Response: resCh}

	toDeviceCh <- req

	time.Sleep(10000 * time.Millisecond)

	select {
	case m := <-resCh:
		if m[0] == 0 {
			t.Error("Response to requests indicate it was not fullfilled")
		}
	default:
		t.Error("No response on activate freestyle message request")
	}
}

func TestActivateGps(t *testing.T) {

	time.Sleep(10000 * time.Millisecond)

	resCh := make(chan []byte)

	buf := []byte{}
	buf = append(buf, []byte(phoneNumber)...)
	buf = append(buf, 0x1B)
	buf = append(buf, []byte(pin)...)
	buf = append(buf, 0x1B)

	req := &CustomProtocol.Request{Id: CustomProtocol.AssignRequestId(), Destination: CustomProtocol.DeviceGPS, Source: CustomProtocol.Web,
		OpCode: CustomProtocol.ActivateGPS, Payload: buf, Response: resCh}

	toDeviceCh <- req

	time.Sleep(10000 * time.Millisecond)

	select {
	case m := <-resCh:
		if m[0] == 0 {
			t.Error("Response to requests indicate it was not fullfilled")
		}
	default:
		t.Error("No response on activate gps request")
	}
}

func TestActivateGeofence(t *testing.T) {

	resCh := make(chan []byte)

	buf := []byte{}
	buf = append(buf, []byte(phoneNumber)...)
	buf = append(buf, 0x1B)
	buf = append(buf, []byte(pin)...)
	buf = append(buf, 0x1B)
	buf = append(buf, []byte("1")...)
	buf = append(buf, 0x1B)
	buf = append(buf, []byte("250")...)
	buf = append(buf, 0x1B)

	req := &CustomProtocol.Request{Id: CustomProtocol.AssignRequestId(), Destination: CustomProtocol.DeviceGPS, Source: CustomProtocol.Web,
		OpCode: CustomProtocol.ActivateGeofence, Payload: buf, Response: resCh}

	toDeviceCh <- req

	time.Sleep(10000 * time.Millisecond)

	select {
	case m := <-resCh:
		if m[0] == 0 {
			t.Error("Response to requests indicate it was not fullfilled")
		}
	default:
		t.Error("No response on activate geofence request")
	}
}

func TestSleepGeogram(t *testing.T) {

	resCh := make(chan []byte)

	buf := []byte{}
	buf = append(buf, []byte(phoneNumber)...)
	buf = append(buf, 0x1B)
	buf = append(buf, []byte(pin)...)
	buf = append(buf, 0x1B)

	req := &CustomProtocol.Request{Id: CustomProtocol.AssignRequestId(), Destination: CustomProtocol.DeviceGPS, Source: CustomProtocol.Web,
		OpCode: CustomProtocol.SleepGeogram, Payload: buf, Response: resCh}

	toDeviceCh <- req

	time.Sleep(10000 * time.Millisecond)

	select {
	case m := <-resCh:
		if m[0] == 0 {
			t.Error("Response to requests indicate it was not fullfilled")
		}
	default:
		t.Error("No response on activate sleep geogram request")
	}
}

func TestActivateIntervalGps(t *testing.T) {

	resCh := make(chan []byte)

	buf := []byte{}
	buf = append(buf, []byte(phoneNumber)...)
	buf = append(buf, 0x1B)
	buf = append(buf, []byte(pin)...)
	buf = append(buf, 0x1B)
	buf = append(buf, []byte("30")...)
	buf = append(buf, 0x1B)

	req := &CustomProtocol.Request{Id: CustomProtocol.AssignRequestId(), Destination: CustomProtocol.DeviceGPS, Source: CustomProtocol.Web,
		OpCode: CustomProtocol.ActivateIntervalGps, Payload: buf, Response: resCh}

	toDeviceCh <- req

	time.Sleep(10000 * time.Millisecond)

	select {
	case m := <-resCh:
		if m[0] == 0 {
			t.Error("Response to requests indicate it was not fullfilled")
		}
	default:
		t.Error("No response on activate interval gps request")
	}
}

func TestSetGeofence(t *testing.T) {

	resCh := make(chan []byte)

	buf := []byte{}
	buf = append(buf, []byte(phoneNumber)...)
	buf = append(buf, 0x1B)
	buf = append(buf, []byte(pin)...)
	buf = append(buf, 0x1B)
	buf = append(buf, []byte("12345678")...)
	buf = append(buf, 0x1B)
	buf = append(buf, []byte("098765432")...)
	buf = append(buf, 0x1B)

	req := &CustomProtocol.Request{Id: CustomProtocol.AssignRequestId(), Destination: CustomProtocol.DeviceGPS, Source: CustomProtocol.Web,
		OpCode: CustomProtocol.SetGeofence, Payload: buf, Response: resCh}

	toDeviceCh <- req

	time.Sleep(20000 * time.Millisecond)

	select {
	case m := <-resCh:
		if m[0] == 0 {
			t.Error("Response to requests indicate it was not fullfilled")
		}
	default:
		t.Error("No response on activate set geofence request")
	}
}

func TestGeogramSetup(t *testing.T) {

	resCh := make(chan []byte)

	buf := []byte{}
	buf = append(buf, []byte(phoneNumber)...)
	buf = append(buf, 0x1B)
	buf = append(buf, []byte(pin)...)
	buf = append(buf, 0x1B)

	req := &CustomProtocol.Request{Id: CustomProtocol.AssignRequestId(), Destination: CustomProtocol.DeviceGPS, Source: CustomProtocol.Web,
		OpCode: CustomProtocol.GeogramSetup, Payload: buf, Response: resCh}

	toDeviceCh <- req

	time.Sleep(30000 * time.Millisecond)

	select {
	case m := <-resCh:
		if m[0] == 0 {
			t.Error("Response to requests indicate it was not fullfilled")
		}
	default:
		t.Error("No response on activate set geofence request")
	}
}
