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
type DeleteLink struct {
	BaseMsg
	res1    byte
	Reply   bool
	All     bool
	Orderly bool
	res2    byte
	Link    uint8
	RsnCode DelLinkRsnCode
	res3    [35]byte
}

// Parse fills the deleteLink fields from the LLC delete link message in buffer
func (d *DeleteLink) Parse(buffer []byte) {
	// init base message fields
	d.SetBaseMsg(buffer)
	buffer = buffer[2:]

	// Reserved 1 byte
	d.res1 = buffer[0]
	buffer = buffer[1:]

	// Reply is first bit in this byte
	d.Reply = (buffer[0] & 0b10000000) > 0

	// All is the next bit in this byte
	d.All = (buffer[0] & 0b01000000) > 0

	// Orderly is the next bit in this byte
	d.Orderly = (buffer[0] & 0b00100000) > 0

	// Remainder of this byte is reserved
	d.res2 = buffer[0] & 0b00011111
	buffer = buffer[1:]

	// Link is 1 byte
	d.Link = buffer[0]
	buffer = buffer[1:]

	// Reason Code is 4 bytes
	d.RsnCode = DelLinkRsnCode(binary.BigEndian.Uint32(buffer[0:4]))
	buffer = buffer[4:]

	// Rest of message is reserved
	copy(d.res3[:], buffer[:])
}

// String converts the delete link message to a string
func (d *DeleteLink) String() string {
	dFmt := "LLC Delete Link: Type: %d, Length: %d, Reply: %t, All: %t, " +
		"Orderly: %t, Link: %d, Reason Code: %s\n"
	return fmt.Sprintf(dFmt, d.Type, d.Length, d.Reply, d.All, d.Orderly,
		d.Link, d.RsnCode)
}

// Reserved converts the delete link message to a string including reserved
// fields
func (d *DeleteLink) Reserved() string {
	dFmt := "LLC Delete Link: Type: %d, Length: %d, Reserved: %#x, " +
		"Reply: %t, All: %t, Orderly: %t, Reserved: %#x, Link: %d, " +
		"Reason Code: %s, Reserved: %#x\n"
	return fmt.Sprintf(dFmt, d.Type, d.Length, d.res1, d.Reply, d.All,
		d.Orderly, d.res2, d.Link, d.RsnCode, d.res3)
}

// ParseDeleteLink parses the LLC delete link message in buffer
func ParseDeleteLink(buffer []byte) *DeleteLink {
	var del DeleteLink
	del.Parse(buffer)
	return &del
}
