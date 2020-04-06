package llc

import "encoding/hex"

// baseMsg stores common message fields
type baseMsg struct {
	raw    []byte
	typ    int
	length int
}

// setRaw stores raw message bytes in the message
func (b *baseMsg) setRaw(buffer []byte) {
	b.raw = make([]byte, len(buffer))
	copy(b.raw[:], buffer[:])
}

// setBaseMsg initializes base message from buffer
func (b *baseMsg) setBaseMsg(buffer []byte) {
	// save raw message bytes
	b.setRaw(buffer)

	// Message type is 1 byte
	b.typ = int(buffer[0])
	buffer = buffer[1:]

	// Message length is 1 byte, should be equal to 44
	b.length = int(buffer[0])
}

// Hex converts the message to a hex dump string
func (b *baseMsg) Hex() string {
	return hex.Dump(b.raw)
}

// GetType returns the type of the message
func (b *baseMsg) GetType() int {
	return b.typ
}
