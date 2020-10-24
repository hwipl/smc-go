package clc

import "bytes"

// SMC-R and SMC-D eyecatchers
var (
	SMCREyecatcher = []byte{0xE2, 0xD4, 0xC3, 0xD9}
	SMCDEyecatcher = []byte{0xE2, 0xD4, 0xC3, 0xC4}
)

const (
	// EyecatcherLen is the length of an eyecatcher
	EyecatcherLen = 4
)

// Eyecatcher stores a SMC Eyecatcher
type Eyecatcher [EyecatcherLen]byte

// String converts the eyecatcher to a string
func (e Eyecatcher) String() string {
	if bytes.Compare(e[:], SMCREyecatcher) == 0 {
		return "SMC-R"
	}
	if bytes.Compare(e[:], SMCDEyecatcher) == 0 {
		return "SMC-D"
	}
	return "Unknown"
}

// HasEyecatcher checks if there is a SMC-R or SMC-D eyecatcher in buf
func HasEyecatcher(buf []byte) bool {
	if bytes.Compare(buf[:EyecatcherLen], SMCREyecatcher) == 0 {
		return true
	}
	if bytes.Compare(buf[:EyecatcherLen], SMCDEyecatcher) == 0 {
		return true
	}
	return false
}
