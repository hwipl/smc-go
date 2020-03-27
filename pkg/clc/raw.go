package clc

import "encoding/hex"

// raw stores the raw bytes of a CLC message
type raw []byte

// Parse saves buf as raw message bytes
func (r *raw) Parse(buf []byte) {
	*r = buf
}

// Dump returns the raw bytes of the message as hex dump string
func (r raw) Dump() string {
	return hex.Dump(r)
}
