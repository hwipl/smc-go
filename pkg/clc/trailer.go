package clc

import "log"

const (
	trailerLen = EyecatcherLen
)

// trailer stores a CLC message trailer
type trailer Eyecatcher

// Parse parses the CLC message trailer at the end of buf
func (t *trailer) Parse(buf []byte) {
	copy(t[:], buf[len(buf)-trailerLen:])
	if !HasEyecatcher(t[:]) {
		log.Println("Error parsing CLC message: invalid trailer")
		errDump(buf[len(buf)-trailerLen:])
		return
	}
}

// String converts the message trailer to a string
func (t trailer) String() string {
	return Eyecatcher(t).String()
}
