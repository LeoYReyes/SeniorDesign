//TODO: Refactor: processRequest that all components can share
package CustomProtocol

import (
	"time"
)

const (
	// Op Code format	ddd fffff
	//	ddd = destination
	//	fffff = function

	// Broadcast opcodes 0 - 32

	// Database opcodes 32 - 63
	NewAccount             = 32
	NewDevice              = 33
	UpdateDeviceGPS        = 34
	UpdateDeviceIP         = 35
	UpdateDeviceKeylog     = 36
	VerifyLoginCredentials = 37
	SetAccount             = 38
	GetDevice              = 39
	SetDevice              = 40
	GetDeviceList          = 41

	// Web opcodes 64 - 95
	UpdateWebMap = 64

	// Device opcodes 96 - 159

	// Geogram opcodes 96 - 127
	ActivateGPS         = 96
	SleepGeogram        = 97
	ActivateGeofence    = 98
	ActivateIntervalGps = 99
	SetGeofence         = 100
	GeogramSetup        = 101
	FreestyleMsg        = 102

	// Laptop opcodes 128 - 159
	CheckDeviceStolen     = 128
	GetUserData           = 129
	UpdateUserKeylogData  = 130
	UpdateUserIPTraceData = 131
	FlagStolen            = 132
	FlagNotStolen         = 133

	// No OpCode 255
	NoOp = 255
)

// Destination constants
const (
	Broadcast    = 0
	Database     = 1
	Web          = 2
	DeviceGPS    = 3
	DeviceLaptop = 4
)

var RequestId = 0

/*
	Destination / Sources

	Broadcast 	== 0
	Database 	== 1
	Web			== 2
	Device		== 3
*/
type Request struct {
	// Unique id
	Id          int
	Destination int
	Source      int
	OpCode      byte
	Payload     []byte
	Response    chan []byte
}

func AssignRequestId() int {
	RequestId += 1
	return RequestId
}

func (req *Request) isThisForMe(i int) bool {
	return true
}

func ParsePayload(payload []byte) []string {
	str := []string{}
	pos := 0
	for index, element := range payload {
		if element == 0x1B {
			str = append(str, string(payload[pos:index]))
			pos = index + 1
		}
	}
	return str
}

/**
 * Creates a payload of strings converted to bytes
 * and separated by escape chars
 */
func CreatePayload(args ...string) []byte {
	var payload []byte
	for _, str := range args {
		payload = append(payload, []byte(str)...)
		payload = append(payload, 0x1B)
	}
	return payload
}

/**
 * Listen to a channel with a set timeout. If there is a response in that
 * time, the function returns true, with the byte array that was sent
 * to the channel. If it does not receive a response before the timeout,
 * then it returns false with a nil value for the array. The timeout is in
 * seconds.
 */
func GetResponse(respCh chan []byte, timeout time.Duration) (bool, []byte) {
	select {
	case response := <-respCh:
		return true, response
	case <-time.After(timeout * time.Second):
		return false, nil
	}
}
