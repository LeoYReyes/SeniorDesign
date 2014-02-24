package CustomRequest

type opCode byte

const (
	checkIn

	// Web opcodes
	updateWebMap

	// Geogram opcodes
	activateGPS

	// Laptop opcodes
	deviceStolen
	deviceNotStolen
	activateKeylog
	activateTraceRoute

	// Database opcodes
	newAccount
	newDevice
	updateDeviceGPS
	updateDeviceIP
	updateDeviceKeylog
	getAccount
	setAccount
	getDevice
	setDevice
	getDeviceList
)

/*
	Destination / Sources

	Broadcast 	== 0
	Database 	== 1
	Web			== 2
	Device		== 3
*/
type CustomRequest struct {
	id          int
	destination int
	source      int
	opCode      byte
	payload     string
}
