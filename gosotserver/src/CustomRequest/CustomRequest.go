package CustomRequest

type opCode byte

const (
	CheckIn = iota

	// Web opcodes
	UpdateWebMap

	// Geogram opcodes
	ActivateGPS

	// Laptop opcodes
	DeviceStolen
	DeviceNotStolen
	ActivateKeylog
	ActivateTraceRoute

	// Database opcodes
	NewAccount
	NewDevice
	UpdateDeviceGPS
	UpdateDeviceIP
	UpdateDeviceKeylog
	GetAccount
	SetAccount
	GetDevice
	SetDevice
	GetDeviceList
)

/*
	Destination / Sources

	Broadcast 	== 0
	Database 	== 1
	Web			== 2
	Device		== 3
*/
type Request struct {
	Id          int
	Destination int
	Source      int
	OpCode      byte
	Payload     string
}
