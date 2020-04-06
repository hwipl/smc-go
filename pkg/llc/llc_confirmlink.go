package llc

import (
	"encoding/binary"
	"fmt"
	"net"
)

// ConfirmLink stores a LLC confirm message
type ConfirmLink struct {
	BaseMsg
	res1             byte
	Reply            bool
	res2             byte
	SenderMAC        net.HardwareAddr
	SenderGID        net.IP
	SenderQP         uint32
	Link             uint8
	SenderLinkUserID uint32
	MaxLinks         uint8
	res3             [9]byte
}

// Parse fills the confirmLink fields from the LLC confirm link message in
// buffer
func (c *ConfirmLink) Parse(buffer []byte) {
	// init base message fields
	c.SetBaseMsg(buffer)
	buffer = buffer[2:]

	// Reserved 1 byte
	c.res1 = buffer[0]
	buffer = buffer[1:]

	// Reply is first bit in this byte
	c.Reply = (buffer[0] & 0b10000000) > 0

	// Remainder of this byte is reserved
	c.res2 = buffer[0] & 0b01111111
	buffer = buffer[1:]

	// sender MAC is a 6 byte MAC address
	c.SenderMAC = make(net.HardwareAddr, 6)
	copy(c.SenderMAC[:], buffer[0:6])
	buffer = buffer[6:]

	// sender GID is an 16 bytes IPv6 address
	c.SenderGID = make(net.IP, net.IPv6len)
	copy(c.SenderGID[:], buffer[0:16])
	buffer = buffer[16:]

	// QP number is 3 bytes
	c.SenderQP = uint32(buffer[0]) << 16
	c.SenderQP |= uint32(buffer[1]) << 8
	c.SenderQP |= uint32(buffer[2])
	buffer = buffer[3:]

	// Link is 1 byte
	c.Link = buffer[0]
	buffer = buffer[1:]

	// Link User ID is 4 bytes
	c.SenderLinkUserID = binary.BigEndian.Uint32(buffer[0:4])
	buffer = buffer[4:]

	// Max Links is 1 byte
	c.MaxLinks = buffer[0]
	buffer = buffer[1:]

	// Rest of message is reserved
	copy(c.res3[:], buffer[:])
}

// String converts the LLC confirm link message to string
func (c *ConfirmLink) String() string {
	cFmt := "LLC Confirm Link: Type: %d, Length: %d, Reply: %t, " +
		"Sender MAC: %s, Sender GID: %s, Sender QP: %d, Link: %d, " +
		"Sender Link UserID: %d, Max Links: %d\n"
	return fmt.Sprintf(cFmt, c.Type, c.Length, c.Reply, c.SenderMAC,
		c.SenderGID, c.SenderQP, c.Link, c.SenderLinkUserID,
		c.MaxLinks)
}

// Reserved converts the LLC confirm link message to string including reserved
// fields
func (c *ConfirmLink) Reserved() string {
	cFmt := "LLC Confirm Link: Type: %d, Length: %d, Reserved: %#x, " +
		"Reply: %t, Reserved: %#x, Sender MAC: %s, Sender GID: %s, " +
		"Sender QP: %d, Link: %d, Sender Link UserID: %d, " +
		"Max Links: %d, Reserved: %#x\n"
	return fmt.Sprintf(cFmt, c.Type, c.Length, c.res1, c.Reply, c.res2,
		c.SenderMAC, c.SenderGID, c.SenderQP, c.Link,
		c.SenderLinkUserID, c.MaxLinks, c.res3)
}

// ParseConfirm parses and prints the LLC confirm link message in buffer
func ParseConfirm(buffer []byte) *ConfirmLink {
	var confirm ConfirmLink
	confirm.Parse(buffer)
	return &confirm
}
