package CustomProtocol

import ()

type Response struct {
	CustomProtocol
	// Id should be the Id of the Request this is a Response to.
	Id          int
	Destination int
	Source      int
	Payload     []byte
}
