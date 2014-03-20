package device

import (
	"CustomProtocol"
	"testing"
	"time"
)

// use a phone number to manually test if the output was correct
var phoneNumber = "2565414217"
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

func TestActivateGps(t *testing.T) {
	//toDeviceCh := make(chan *CustomProtocol.Request)
	//fromDeviceCh := make(chan *CustomProtocol.Request)
	go StartDeviceServer(fromDeviceCh, toDeviceCh)
	go hack()

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
	//toDeviceCh := make(chan *CustomProtocol.Request)
	//fromDeviceCh := make(chan *CustomProtocol.Request)
	//go StartDeviceServer(fromDeviceCh, toDeviceCh)

	//time.Sleep(10000 * time.Millisecond)

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
		t.Error("No response on activate gps request")
	}
}

func TestSleepGeogram(t *testing.T) {
	//toDeviceCh := make(chan *CustomProtocol.Request)
	//fromDeviceCh := make(chan *CustomProtocol.Request)
	//go StartDeviceServer(fromDeviceCh, toDeviceCh)

	//time.Sleep(10000 * time.Millisecond)

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
		t.Error("No response on activate gps request")
	}
}

func TestActivateIntervalGps(t *testing.T) {
	//toDeviceCh := make(chan *CustomProtocol.Request)
	//fromDeviceCh := make(chan *CustomProtocol.Request)
	//go StartDeviceServer(fromDeviceCh, toDeviceCh)

	//time.Sleep(10000 * time.Millisecond)

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
		t.Error("No response on activate gps request")
	}
}

func TestSetGeofence(t *testing.T) {
	//toDeviceCh := make(chan *CustomProtocol.Request)
	//fromDeviceCh := make(chan *CustomProtocol.Request)
	//go StartDeviceServer(fromDeviceCh, toDeviceCh)

	//time.Sleep(10000 * time.Millisecond)

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
