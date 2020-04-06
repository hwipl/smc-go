package llc

import (
	"encoding/binary"
	"fmt"
)

// DelLinkRsnCode stores the reason code of a delete link message
type DelLinkRsnCode uint32

// String converts the delete link reason code to a string
func (d DelLinkRsnCode) String() string {
	var rsn string

	switch d {
	case 0x00010000:
		rsn = "Lost path"

	case 0x00020000:
		rsn = "Operator initiated termination"

	case 0x00030000:
		rsn = "Program initiated termination (link inactivity)"

	case 0x00040000:
		rsn = "LLC protocol violation"

	case 0x00050000:
		rsn = "Asymmetric link no longer needed"
	case 0x00100000:
		rsn = "Unknown link ID (no link)"
	default:
		rsn = "unknown"
	}

	return fmt.Sprintf("%d (%s)", d, rsn)
}

// delteLink stores a LLC delete link message
type deleteLink struct {
	BaseMsg
	res1    byte
	reply   bool
	all     bool
	orderly bool
	res2    byte
	link    uint8
	rsnCode DelLinkRsnCode
	res3    [35]byte
}

// Parse fills the deleteLink fields from the LLC delete link message in buffer
func (d *deleteLink) Parse(buffer []byte) {
	// init base message fields
	d.SetBaseMsg(buffer)
	buffer = buffer[2:]

	// Reserved 1 byte
	d.res1 = buffer[0]
	buffer = buffer[1:]

	// Reply is first bit in this byte
	d.reply = (buffer[0] & 0b10000000) > 0

	// All is the next bit in this byte
	d.all = (buffer[0] & 0b01000000) > 0

	// Orderly is the next bit in this byte
	d.orderly = (buffer[0] & 0b00100000) > 0

	// Remainder of this byte is reserved
	d.res2 = buffer[0] & 0b00011111
	buffer = buffer[1:]

	// Link is 1 byte
	d.link = buffer[0]
	buffer = buffer[1:]

	// Reason Code is 4 bytes
	d.rsnCode = DelLinkRsnCode(binary.BigEndian.Uint32(buffer[0:4]))
	buffer = buffer[4:]

	// Rest of message is reserved
	copy(d.res3[:], buffer[:])
}

// String converts the delete link message to a string
func (d *deleteLink) String() string {
	dFmt := "LLC Delete Link: Type: %d, Length: %d, Reply: %t, All: %t, " +
		"Orderly: %t, Link: %d, Reason Code: %s\n"
	return fmt.Sprintf(dFmt, d.Type, d.Length, d.reply, d.all, d.orderly,
		d.link, d.rsnCode)
}

// Reserved converts the delete link message to a string including reserved
// fields
func (d *deleteLink) Reserved() string {
	dFmt := "LLC Delete Link: Type: %d, Length: %d, Reserved: %#x, " +
		"Reply: %t, All: %t, Orderly: %t, Reserved: %#x, Link: %d, " +
		"Reason Code: %s, Reserved: %#x\n"
	return fmt.Sprintf(dFmt, d.Type, d.Length, d.res1, d.reply, d.all,
		d.orderly, d.res2, d.link, d.rsnCode, d.res3)
}

// parseDeleteLink parses the LLC delete link message in buffer
func parseDeleteLink(buffer []byte) *deleteLink {
	var del deleteLink
	del.Parse(buffer)
	return &del
}
