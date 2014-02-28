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
}

func AssignRequestId() int {
	RequestId += 1
	return RequestId
}

func (req *Request) isThisForMe(i int) bool {
	return true
}

/*//tests if request is for the calling source
func (req *Request) isThisForMe(int i) bool {
	if req.Destination = i  {
		return true;
	} else {
		return false;
	}
}

//Check OpCode and Destination are Valid.
func (req *Request) checkValidity() bool{
	req.opCodeExists()
}


//Parses Payload and returns all string variables
//If wants it to not be a string can change the type- str2byteExample: []byte("string")
func (req *Request) parsePayload()  {
	return strings.Split(req.Payload, "/")
}


//parses Payload with delimter, but max is the number n.
//Might help with error checking.
func (req *Request) parsePayload(int n)  {
	return strings.SplitN(req.Payload, "/", n)
}


// Checks for valid OpCode and Sends to check exists in Destination.
func (req *Request) opCodeExists() {
	opcode =  req.OpCode
	switch opcode {
		default:ErrorOpCodeDNE()
		case #, #, #, #, : isFunctionInDestination(0) //Broadcast Destination.
		case #, #, #, #, : isFunctionInDestination(1) //Database Destination..
		case #, #, #, #, : isFunctionInDestination(2) //Web Destination.
		case #, #, #, #, : isFunctionInDestination(3) //Device Destination.
	}
}


//Test OpCode Value against Destination Function Values.
//VALUES NOT SET RIGHT YET FOR NESTED IF STATEMENT

func (req *Request) isFunctionInDestination(int i) bool {
	//Destination DNE
	if req.Destination >= 0 && req.Destination <= 3
		Destination := string(req.Destination)
		fmt.Printf("%q\n", "Destination '", Destination, "' does not exist.")
		return false
	//Destination and OpCode Match
	} else if i = req.Destination  {
		return
	//Destination and OpCode Mismatch
	} else {
		Destination := string(req.Destination)
		OpCode := string(req.OpCode)
		fmt.Printf("%q\n\t", "OpCode function called doesn't belong to that Destination.")
		fmt.Printf("%q\n\t", "Destination: ", Destination)
		fmt.Printf("%q\n", "OpCode: ", OpCode)
		return false
	}
}


func (req *Request) ErrorOpCodeDNE() bool{
	OpCode := string(req.OpCode)
	fmt.Printf("%q\n", "OpCode '", OpCode, "' does not exist.")
	return false
}*/
