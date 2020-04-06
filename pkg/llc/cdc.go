package llc

import (
	"encoding/binary"
	"fmt"
)

// cdc stores a CDC message
type cdc struct {
	baseMsg
	seqNum   uint16
	alertTkn uint32
	res1     [2]byte
	prodWrap uint16
	prodCurs uint32
	res2     [2]byte
	consWrap uint16
	consCurs uint32
	b        bool
	p        bool
	u        bool
	r        bool
	f        bool
	res3     byte
	d        bool
	c        bool
	a        bool
	res4     [19]byte
}

// Parse fills the cdc fields from the CDC message in buffer
func (c *cdc) Parse(buffer []byte) {
	// save raw message bytes
	c.setRaw(buffer)

	// Message type is 1 byte
	c.Type = int(buffer[0])
	buffer = buffer[1:]

	// Message length is 1 byte, should be equal to 44
	c.Length = int(buffer[0])
	buffer = buffer[1:]

	// Sequence number is 2 bytes
	c.seqNum = binary.BigEndian.Uint16(buffer[0:2])
	buffer = buffer[2:]

	// Alert token is 4 bytes
	c.alertTkn = binary.BigEndian.Uint32(buffer[0:4])
	buffer = buffer[4:]

	// Reserved are 2 bytes
	copy(c.res1[:], buffer[0:2])
	buffer = buffer[2:]

	// Producer cursor wrap sequence number is 2 bytes
	c.prodWrap = binary.BigEndian.Uint16(buffer[0:2])
	buffer = buffer[2:]

	// Producer cursor is 4 bytes
	c.prodCurs = binary.BigEndian.Uint32(buffer[0:4])
	buffer = buffer[4:]

	// Reserved are 2 bytes
	copy(c.res2[:], buffer[0:2])
	buffer = buffer[2:]

	// Consumer cursor wrap sequence number is 2 bytes
	c.consWrap = binary.BigEndian.Uint16(buffer[0:2])
	buffer = buffer[2:]

	// Consumer cursor is 4 bytes
	c.consCurs = binary.BigEndian.Uint32(buffer[0:4])
	buffer = buffer[4:]

	// B-bit/Writer blocked indicator is the first bit in this byte
	c.b = (buffer[0] & 0b10000000) > 0

	// P-bit/Urgent data pending is next bit in this byte
	c.p = (buffer[0] & 0b01000000) > 0

	// U-bit/Urgent data present is next bit in this byte
	c.u = (buffer[0] & 0b00100000) > 0

	// R-bit/Request for consumer cursor update is next bit in this byte
	c.r = (buffer[0] & 0b00010000) > 0

	// F-bit/Failover validation indicator is next bit in this byte
	c.f = (buffer[0] & 0b00001000) > 0

	// Reserved are the remaining bits in this byte
	c.res3 = buffer[0] & 0b00000111
	buffer = buffer[1:]

	// D-bit/Sending done indicator is the first bit in this byte
	c.d = (buffer[0] & 0b10000000) > 0

	// C-bit/PeerConnectionClosed indicator is the next bit in this byte
	c.c = (buffer[0] & 0b01000000) > 0

	// A-bit/Abnormal close indicator is the next bit in this byte
	c.a = (buffer[0] & 0b00100000) > 0

	// Reserved are the remaining bits in this byte
	c.res4[0] = buffer[0] & 0b00011111
	buffer = buffer[1:]

	// Rest of message is reserved
	copy(c.res4[1:], buffer[:])
}

// String converts the cdc message into a string
func (c *cdc) String() string {
	cFmt := "CDC: Type: %d, Length %d, Sequence Number: %d, " +
		"Alert Token: %d, Producer Wrap: %d, " +
		"Producer Cursor: %d, Consumer Wrap: %d, " +
		"Consumer Cursor: %d, Writer Blocked: %t, " +
		"Urgent Data Pending: %t, Urgent Data Present: %t, " +
		"Request for Consumer Cursor Update: %t, " +
		"Failover Validation: %t, Sending Done: %t, " +
		"Peer Connection Closed: %t, Abnormal Close: %t\n"
	return fmt.Sprintf(cFmt, c.Type, c.Length, c.seqNum, c.alertTkn,
		c.prodWrap, c.prodCurs, c.consWrap, c.consCurs, c.b, c.p, c.u,
		c.r, c.f, c.d, c.c, c.a)
}

// Reserved converts the cdc message into a string including reserved fields
func (c *cdc) Reserved() string {
	cFmt := "CDC: Type: %d, Length %d, Sequence Number: %d, " +
		"Alert Token: %d, Reserved: %#x, Producer Wrap: %d, " +
		"Producer Cursor: %d, Reserved: %#x, Consumer Wrap: %d, " +
		"Consumer Cursor: %d, Writer Blocked: %t, " +
		"Urgent Data Pending: %t, Urgent Data Present: %t, " +
		"Request for Consumer Cursor Update: %t, " +
		"Failover Validation: %t, Reserved: %#x, Sending Done: %t, " +
		"Peer Connection Closed: %t, Abnormal Close: %t, " +
		"Reserved: %#x\n"
	return fmt.Sprintf(cFmt, c.Type, c.Length, c.seqNum, c.alertTkn, c.res1,
		c.prodWrap, c.prodCurs, c.res2, c.consWrap, c.consCurs, c.b,
		c.p, c.u, c.r, c.f, c.res3, c.d, c.c, c.a, c.res4)
}

// parseCDC parses the CDC message in buffer
func parseCDC(buffer []byte) *cdc {
	var c cdc
	c.Parse(buffer)
	return &c
}
