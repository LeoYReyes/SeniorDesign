package CustomRequest

type opCode byte

const (
	CheckIn = iota

	// Web opcodes
	UpdateWebMap // 1

	// Geogram opcodes
	ActivateGPS // 2

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

	// No OpCode
	noOp
)

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
	Payload     string
}

func (req *Request) isThisForMe(i int) bool {
	return true
}
