package llc

import (
	"encoding/binary"
	"fmt"
	"net"
)

// confirmLink stores a LLC confirm message
type confirmLink struct {
	baseMsg
	res1             byte
	reply            bool
	res2             byte
	senderMAC        net.HardwareAddr
	senderGID        net.IP
	senderQP         uint32
	link             uint8
	senderLinkUserID uint32
	maxLinks         uint8
	res3             [9]byte
}

// Parse fills the confirmLink fields from the LLC confirm link message in
// buffer
func (c *confirmLink) Parse(buffer []byte) {
	// init base message fields
	c.setBaseMsg(buffer)
	buffer = buffer[2:]

	// Reserved 1 byte
	c.res1 = buffer[0]
	buffer = buffer[1:]

	// Reply is first bit in this byte
	c.reply = (buffer[0] & 0b10000000) > 0

	// Remainder of this byte is reserved
	c.res2 = buffer[0] & 0b01111111
	buffer = buffer[1:]

	// sender MAC is a 6 byte MAC address
	c.senderMAC = make(net.HardwareAddr, 6)
	copy(c.senderMAC[:], buffer[0:6])
	buffer = buffer[6:]

	// sender GID is an 16 bytes IPv6 address
	c.senderGID = make(net.IP, net.IPv6len)
	copy(c.senderGID[:], buffer[0:16])
	buffer = buffer[16:]

	// QP number is 3 bytes
	c.senderQP = uint32(buffer[0]) << 16
	c.senderQP |= uint32(buffer[1]) << 8
	c.senderQP |= uint32(buffer[2])
	buffer = buffer[3:]

	// Link is 1 byte
	c.link = buffer[0]
	buffer = buffer[1:]

	// Link User ID is 4 bytes
	c.senderLinkUserID = binary.BigEndian.Uint32(buffer[0:4])
	buffer = buffer[4:]

	// Max Links is 1 byte
	c.maxLinks = buffer[0]
	buffer = buffer[1:]

	// Rest of message is reserved
	copy(c.res3[:], buffer[:])
}

// String converts the LLC confirm link message to string
func (c *confirmLink) String() string {
	cFmt := "LLC Confirm Link: Type: %d, Length: %d, Reply: %t, " +
		"Sender MAC: %s, Sender GID: %s, Sender QP: %d, Link: %d, " +
		"Sender Link UserID: %d, Max Links: %d\n"
	return fmt.Sprintf(cFmt, c.typ, c.length, c.reply, c.senderMAC,
		c.senderGID, c.senderQP, c.link, c.senderLinkUserID,
		c.maxLinks)
}

// Reserved converts the LLC confirm link message to string including reserved
// fields
func (c *confirmLink) Reserved() string {
	cFmt := "LLC Confirm Link: Type: %d, Length: %d, Reserved: %#x, " +
		"Reply: %t, Reserved: %#x, Sender MAC: %s, Sender GID: %s, " +
		"Sender QP: %d, Link: %d, Sender Link UserID: %d, " +
		"Max Links: %d, Reserved: %#x\n"
	return fmt.Sprintf(cFmt, c.typ, c.length, c.res1, c.reply, c.res2,
		c.senderMAC, c.senderGID, c.senderQP, c.link,
		c.senderLinkUserID, c.maxLinks, c.res3)
}

// parseConfirm parses and prints the LLC confirm link message in buffer
func parseConfirm(buffer []byte) *confirmLink {
	var confirm confirmLink
	confirm.Parse(buffer)
	return &confirm
}
