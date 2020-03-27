package clc

import "bytes"

var (
	SMCREyecatcher = []byte{0xE2, 0xD4, 0xC3, 0xD9}
	SMCDEyecatcher = []byte{0xE2, 0xD4, 0xC3, 0xC4}
)

const (
	eyecatcherLen = 4
)

// eyecatcher stores a SMC eyecatcher
type eyecatcher [eyecatcherLen]byte

// String converts the eyecatcher to a string
func (e eyecatcher) String() string {
	if bytes.Compare(e[:], SMCREyecatcher) == 0 {
		return "SMC-R"
	}
	if bytes.Compare(e[:], SMCDEyecatcher) == 0 {
		return "SMC-D"
	}
	return "Unknown"
}

// hasEyecatcher checks if there is a SMC-R or SMC-D eyecatcher in buf
func hasEyecatcher(buf []byte) bool {
	if bytes.Compare(buf[:eyecatcherLen], SMCREyecatcher) == 0 {
		return true
	}
	if bytes.Compare(buf[:eyecatcherLen], SMCDEyecatcher) == 0 {
		return true
	}
	return false
}
