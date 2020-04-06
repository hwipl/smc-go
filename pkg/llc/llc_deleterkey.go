package llc

import (
	"encoding/binary"
	"fmt"
)

// deleteRKey stores a LLC delete RKey message
type deleteRKey struct {
	BaseMsg
	res1      byte
	reply     bool
	res2      byte
	reject    bool // negative response
	res3      byte
	count     uint8
	errorMask byte
	res4      [2]byte
	rkeys     [8]uint32
	res5      [4]byte
}

// Parse fills the deleteRKey fields from the delete RKey message in buffer
func (d *deleteRKey) Parse(buffer []byte) {
	// init base message fields
	d.SetBaseMsg(buffer)
	buffer = buffer[2:]

	// Reserved 1 byte
	d.res1 = buffer[0]
	buffer = buffer[1:]

	// Reply is first bit in this byte
	d.reply = (buffer[0] & 0b10000000) > 0

	// Reserved is the next bit in this byte
	d.res2 = (buffer[0] & 0b01000000) >> 6

	// Negative response flag is the next bit in this byte
	d.reject = (buffer[0] & 0b00100000) > 0

	// Remainder of this byte is reserved
	d.res3 = buffer[0] & 0b00011111
	buffer = buffer[1:]

	// Count is 1 byte
	d.count = buffer[0]
	buffer = buffer[1:]

	// Error Mask is 1 byte
	d.errorMask = buffer[0]
	buffer = buffer[1:]

	// Reserved are 2 bytes
	copy(d.res4[:], buffer[0:2])
	buffer = buffer[2:]

	// Deleted RKeys are each 4 bytes
	// parse
	// * First deleted RKey
	// * Second deleted RKey
	// * Third deleted RKey
	// * Fourth deleted RKey
	// * Fifth deleted RKey
	// * Sixth deleted RKey
	// * Seventh deleted RKey
	// * Eighth deleted RKey
	for i := range d.rkeys {
		d.rkeys[i] = binary.BigEndian.Uint32(buffer[0:4])
		buffer = buffer[4:]
	}

	// Rest of message is reserved
	copy(d.res5[:], buffer[:])
}

// String converts the delete RKey message to a string
func (d *deleteRKey) String() string {
	var rkeys string

	for i := range d.rkeys {
		rkeys += fmt.Sprintf(", RKey %d: %d", i, d.rkeys[i])
	}

	dFmt := "LLC Delete RKey: Type: %d, Length: %d, " +
		"Reply: %t, Negative Response: %t, Count: %d, " +
		"Error Mask: %#b%s\n"
	return fmt.Sprintf(dFmt, d.Type, d.Length, d.reply, d.reject, d.count,
		d.errorMask, rkeys)
}

// Reserved converts the delete RKey message to a string including reserved
// fields
func (d *deleteRKey) Reserved() string {
	var rkeys string

	for i := range d.rkeys {
		rkeys += fmt.Sprintf("RKey %d: %d, ", i, d.rkeys[i])
	}

	dFmt := "LLC Delete RKey: Type: %d, Length: %d, " +
		"Reserved: %#x, Reply: %t, Reserved: %#x, " +
		"Negative Response: %t, Reserved: %#x, " +
		"Count: %d, Error Mask: %#b, Reserved: %#x, %sReserved: %#x\n"
	return fmt.Sprintf(dFmt, d.Type, d.Length, d.res1, d.reply, d.res2,
		d.reject, d.res3, d.count, d.errorMask, d.res4, rkeys, d.res5)
}

// parseDeleteRKey parses the LLC delete RKey message in buffer
func parseDeleteRKey(buffer []byte) *deleteRKey {
	var del deleteRKey
	del.Parse(buffer)
	return &del
}
