package llc

import (
	"encoding/binary"
	"fmt"
)

// RMBSpec stores another RMB specificiation
type RMBSpec struct {
	Link  uint8
	RKey  uint32
	VAddr uint64
}

// Parse fills the rmbSpec fields from buffer
func (r *RMBSpec) Parse(buffer []byte) {
	// other link rmb specifications are 13 bytes and consist of:
	// * link number (1 byte)
	// * RMB's RKey for the specified link (4 bytes)
	// * RMB's virtual address for the specified link (8 bytes)
	r.Link = buffer[0]
	r.RKey = binary.BigEndian.Uint32(buffer[1:5])
	r.VAddr = binary.BigEndian.Uint64(buffer[5:13])
}

// String converts the rmbSpec to a string
func (r *RMBSpec) String() string {
	rFmt := "[Link: %d, RKey: %d, Virtual Address: %#x]"
	return fmt.Sprintf(rFmt, r.Link, r.RKey, r.VAddr)
}

// ConfirmRKey stores a LLC confirm RKey message
type ConfirmRKey struct {
	BaseMsg
	res1      byte
	Reply     bool
	res2      byte
	Reject    bool // negative response
	Retry     bool // configuration retry
	res3      byte
	NumTkns   uint8
	RKey      uint32
	VAddr     uint64
	OtherRMBs [2]RMBSpec
	res4      byte
}

// Parse fills the confirmRKey fields from the confirm RKey message in buffer
func (c *ConfirmRKey) Parse(buffer []byte) {
	// init base message fields
	c.SetBaseMsg(buffer)
	buffer = buffer[2:]

	// Reserved 1 byte
	c.res1 = buffer[0]
	buffer = buffer[1:]

	// Reply is first bit in this byte
	c.Reply = (buffer[0] & 0b10000000) > 0

	// Reserved is the next bit in this byte
	c.res2 = (buffer[0] & 0b01000000) >> 6

	// Negative response flag is the next bit in this byte
	c.Reject = (buffer[0] & 0b00100000) > 0

	// Configuration Retry is the next bit in this byte
	c.Retry = (buffer[0] & 0b00010000) > 0

	// Remainder of this byte is reserved
	c.res3 = buffer[0] & 0b00001111
	buffer = buffer[1:]

	// Number of tokens is 1 byte
	c.NumTkns = buffer[0]
	buffer = buffer[1:]

	// New RMB RKey for this link is 4 bytes
	c.RKey = binary.BigEndian.Uint32(buffer[0:4])
	buffer = buffer[4:]

	// New RMB virtual address for this link is 8 bytes
	c.VAddr = binary.BigEndian.Uint64(buffer[0:8])
	buffer = buffer[8:]

	// other link rmb specifications are each 13 bytes
	// parse
	// * first other link rmb (can be all zeros)
	// * second other link rmb (can be all zeros)
	for i := range c.OtherRMBs {
		c.OtherRMBs[i].Parse(buffer)
		buffer = buffer[13:]
	}

	// Rest of message is reserved
	c.res4 = buffer[0]
}

// String converts the confirm RKey message to a string
func (c *ConfirmRKey) String() string {
	var others string

	for i := range c.OtherRMBs {
		others += fmt.Sprintf(", Other Link RMB %d: %s", i+1,
			&c.OtherRMBs[i])
	}

	cFmt := "LLC Confirm RKey: Type: %d, Length: %d, Reply: %t, " +
		"Negative Response: %t, Configuration Retry: %t, " +
		"Number of Tokens: %d, This RKey: %d, This VAddr: %#x%s\n"
	return fmt.Sprintf(cFmt, c.Type, c.Length, c.Reply, c.Reject, c.Retry,
		c.NumTkns, c.RKey, c.VAddr, others)
}

// Reserved converts the confirm RKey message to a string including reserved
// fields
func (c *ConfirmRKey) Reserved() string {
	var others string

	for i := range c.OtherRMBs {
		others += fmt.Sprintf("Other Link RMB %d: %s, ", i+1,
			&c.OtherRMBs[i])
	}

	cFmt := "LLC Confirm RKey: Type: %d, Length: %d, Reserved: %#x, " +
		"Reply: %t, Reserved: %#x, Negative Response: %t, " +
		"Configuration Retry: %t, Reserved: %#x, " +
		"Number of Tokens: %d, This RKey: %d, This VAddr: %#x, " +
		"%sReserved: %#x\n"
	return fmt.Sprintf(cFmt, c.Type, c.Length, c.res1, c.Reply, c.res2,
		c.Reject, c.Retry, c.res3, c.NumTkns, c.RKey, c.VAddr, others,
		c.res4)
}

// ParseConfirmRKey parses the LLC confirm RKey message in buffer
func ParseConfirmRKey(buffer []byte) *ConfirmRKey {
	var confirm ConfirmRKey
	confirm.Parse(buffer)
	return &confirm
}
