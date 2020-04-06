package llc

import (
	"encoding/binary"
	"fmt"
)

// CDC stores a CDC message
type CDC struct {
	BaseMsg
	SeqNum   uint16
	AlertTkn uint32
	res1     [2]byte
	ProdWrap uint16
	ProdCurs uint32
	res2     [2]byte
	ConsWrap uint16
	ConsCurs uint32
	B        bool
	P        bool
	U        bool
	R        bool
	F        bool
	res3     byte
	D        bool
	C        bool
	A        bool
	res4     [19]byte
}

// Parse fills the cdc fields from the CDC message in buffer
func (c *CDC) Parse(buffer []byte) {
	// save raw message bytes
	c.setRaw(buffer)

	// Message type is 1 byte
	c.Type = int(buffer[0])
	buffer = buffer[1:]

	// Message length is 1 byte, should be equal to 44
	c.Length = int(buffer[0])
	buffer = buffer[1:]

	// Sequence number is 2 bytes
	c.SeqNum = binary.BigEndian.Uint16(buffer[0:2])
	buffer = buffer[2:]

	// Alert token is 4 bytes
	c.AlertTkn = binary.BigEndian.Uint32(buffer[0:4])
	buffer = buffer[4:]

	// Reserved are 2 bytes
	copy(c.res1[:], buffer[0:2])
	buffer = buffer[2:]

	// Producer cursor wrap sequence number is 2 bytes
	c.ProdWrap = binary.BigEndian.Uint16(buffer[0:2])
	buffer = buffer[2:]

	// Producer cursor is 4 bytes
	c.ProdCurs = binary.BigEndian.Uint32(buffer[0:4])
	buffer = buffer[4:]

	// Reserved are 2 bytes
	copy(c.res2[:], buffer[0:2])
	buffer = buffer[2:]

	// Consumer cursor wrap sequence number is 2 bytes
	c.ConsWrap = binary.BigEndian.Uint16(buffer[0:2])
	buffer = buffer[2:]

	// Consumer cursor is 4 bytes
	c.ConsCurs = binary.BigEndian.Uint32(buffer[0:4])
	buffer = buffer[4:]

	// B-bit/Writer blocked indicator is the first bit in this byte
	c.B = (buffer[0] & 0b10000000) > 0

	// P-bit/Urgent data pending is next bit in this byte
	c.P = (buffer[0] & 0b01000000) > 0

	// U-bit/Urgent data present is next bit in this byte
	c.U = (buffer[0] & 0b00100000) > 0

	// R-bit/Request for consumer cursor update is next bit in this byte
	c.R = (buffer[0] & 0b00010000) > 0

	// F-bit/Failover validation indicator is next bit in this byte
	c.F = (buffer[0] & 0b00001000) > 0

	// Reserved are the remaining bits in this byte
	c.res3 = buffer[0] & 0b00000111
	buffer = buffer[1:]

	// D-bit/Sending done indicator is the first bit in this byte
	c.D = (buffer[0] & 0b10000000) > 0

	// C-bit/PeerConnectionClosed indicator is the next bit in this byte
	c.C = (buffer[0] & 0b01000000) > 0

	// A-bit/Abnormal close indicator is the next bit in this byte
	c.A = (buffer[0] & 0b00100000) > 0

	// Reserved are the remaining bits in this byte
	c.res4[0] = buffer[0] & 0b00011111
	buffer = buffer[1:]

	// Rest of message is reserved
	copy(c.res4[1:], buffer[:])
}

// String converts the cdc message into a string
func (c *CDC) String() string {
	cFmt := "CDC: Type: %d, Length %d, Sequence Number: %d, " +
		"Alert Token: %d, Producer Wrap: %d, " +
		"Producer Cursor: %d, Consumer Wrap: %d, " +
		"Consumer Cursor: %d, Writer Blocked: %t, " +
		"Urgent Data Pending: %t, Urgent Data Present: %t, " +
		"Request for Consumer Cursor Update: %t, " +
		"Failover Validation: %t, Sending Done: %t, " +
		"Peer Connection Closed: %t, Abnormal Close: %t\n"
	return fmt.Sprintf(cFmt, c.Type, c.Length, c.SeqNum, c.AlertTkn,
		c.ProdWrap, c.ProdCurs, c.ConsWrap, c.ConsCurs, c.B, c.P, c.U,
		c.R, c.F, c.D, c.C, c.A)
}

// Reserved converts the cdc message into a string including reserved fields
func (c *CDC) Reserved() string {
	cFmt := "CDC: Type: %d, Length %d, Sequence Number: %d, " +
		"Alert Token: %d, Reserved: %#x, Producer Wrap: %d, " +
		"Producer Cursor: %d, Reserved: %#x, Consumer Wrap: %d, " +
		"Consumer Cursor: %d, Writer Blocked: %t, " +
		"Urgent Data Pending: %t, Urgent Data Present: %t, " +
		"Request for Consumer Cursor Update: %t, " +
		"Failover Validation: %t, Reserved: %#x, Sending Done: %t, " +
		"Peer Connection Closed: %t, Abnormal Close: %t, " +
		"Reserved: %#x\n"
	return fmt.Sprintf(cFmt, c.Type, c.Length, c.SeqNum, c.AlertTkn, c.res1,
		c.ProdWrap, c.ProdCurs, c.res2, c.ConsWrap, c.ConsCurs, c.B,
		c.P, c.U, c.R, c.F, c.res3, c.D, c.C, c.A, c.res4)
}

// ParseCDC parses the CDC message in buffer
func ParseCDC(buffer []byte) *CDC {
	var c CDC
	c.Parse(buffer)
	return &c
}
