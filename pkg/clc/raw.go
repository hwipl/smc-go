package clc

import "encoding/hex"

// Raw stores the Raw bytes of a CLC message
type Raw []byte

// Parse saves buf as raw message bytes
func (r *Raw) Parse(buf []byte) {
	*r = buf
}

// Dump returns the raw bytes of the message as hex dump string
func (r Raw) Dump() string {
	return hex.Dump(r)
}
