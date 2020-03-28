package clc

import "log"

const (
	TrailerLen = EyecatcherLen
)

// Trailer stores a CLC message Trailer
type Trailer Eyecatcher

// Parse parses the CLC message trailer at the end of buf
func (t *Trailer) Parse(buf []byte) {
	copy(t[:], buf[len(buf)-TrailerLen:])
	if !HasEyecatcher(t[:]) {
		log.Println("Error parsing CLC message: invalid trailer")
		errDump(buf[len(buf)-TrailerLen:])
		return
	}
}

// String converts the message trailer to a string
func (t Trailer) String() string {
	return Eyecatcher(t).String()
}
