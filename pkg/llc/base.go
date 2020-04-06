package llc

import "encoding/hex"

// baseMsg stores common message fields
type baseMsg struct {
	Raw    []byte
	Type   int
	Length int
}

// setRaw stores raw message bytes in the message
func (b *baseMsg) setRaw(buffer []byte) {
	b.Raw = make([]byte, len(buffer))
	copy(b.Raw[:], buffer[:])
}

// SetBaseMsg initializes base message from buffer
func (b *baseMsg) SetBaseMsg(buffer []byte) {
	// save raw message bytes
	b.setRaw(buffer)

	// Message type is 1 byte
	b.Type = int(buffer[0])
	buffer = buffer[1:]

	// Message length is 1 byte, should be equal to 44
	b.Length = int(buffer[0])
}

// Hex converts the message to a hex dump string
func (b *baseMsg) Hex() string {
	return hex.Dump(b.Raw)
}

// GetType returns the type of the message
func (b *baseMsg) GetType() int {
	return b.Type
}
