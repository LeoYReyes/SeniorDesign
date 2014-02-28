package CustomProtocol

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
	ActivateGPS = 96

	// Laptop opcodes 128 - 159
	DeviceStolen       = 128
	DeviceNotStolen    = 129
	ActivateKeylog     = 130
	ActivateTraceRoute = 131

	// No OpCode 255
	noOp = 255
)

// Destination constants
const (
	Broadcast    = 0
	Database     = 1
	Web          = 2
	DeviceGPS    = 3
	DeviceLaptop = 4
)

type CustomProtocol struct {
	Id          int
	Destination int
	Source      int
	Payload     []byte
}
