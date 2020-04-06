package llc

import (
	"encoding/binary"
	"fmt"
)

// rmbSpec stores another RMB specificiation
type rmbSpec struct {
	link  uint8
	rkey  uint32
	vaddr uint64
}

// parse fills the rmbSpec fields from buffer
func (r *rmbSpec) parse(buffer []byte) {
	// other link rmb specifications are 13 bytes and consist of:
	// * link number (1 byte)
	// * RMB's RKey for the specified link (4 bytes)
	// * RMB's virtual address for the specified link (8 bytes)
	r.link = buffer[0]
	r.rkey = binary.BigEndian.Uint32(buffer[1:5])
	r.vaddr = binary.BigEndian.Uint64(buffer[5:13])
}

// String converts the rmbSpec to a string
func (r *rmbSpec) String() string {
	rFmt := "[Link: %d, RKey: %d, Virtual Address: %#x]"
	return fmt.Sprintf(rFmt, r.link, r.rkey, r.vaddr)
}

// confirmRKey stores a LLC confirm RKey message
type confirmRKey struct {
	baseMsg
	res1      byte
	reply     bool
	res2      byte
	reject    bool // negative response
	retry     bool // configuration retry
	res3      byte
	numTkns   uint8
	rkey      uint32
	vaddr     uint64
	otherRMBs [2]rmbSpec
	res4      byte
}

// Parse fills the confirmRKey fields from the confirm RKey message in buffer
func (c *confirmRKey) Parse(buffer []byte) {
	// init base message fields
	c.SetBaseMsg(buffer)
	buffer = buffer[2:]

	// Reserved 1 byte
	c.res1 = buffer[0]
	buffer = buffer[1:]

	// Reply is first bit in this byte
	c.reply = (buffer[0] & 0b10000000) > 0

	// Reserved is the next bit in this byte
	c.res2 = (buffer[0] & 0b01000000) >> 6

	// Negative response flag is the next bit in this byte
	c.reject = (buffer[0] & 0b00100000) > 0

	// Configuration Retry is the next bit in this byte
	c.retry = (buffer[0] & 0b00010000) > 0

	// Remainder of this byte is reserved
	c.res3 = buffer[0] & 0b00001111
	buffer = buffer[1:]

	// Number of tokens is 1 byte
	c.numTkns = buffer[0]
	buffer = buffer[1:]

	// New RMB RKey for this link is 4 bytes
	c.rkey = binary.BigEndian.Uint32(buffer[0:4])
	buffer = buffer[4:]

	// New RMB virtual address for this link is 8 bytes
	c.vaddr = binary.BigEndian.Uint64(buffer[0:8])
	buffer = buffer[8:]

	// other link rmb specifications are each 13 bytes
	// parse
	// * first other link rmb (can be all zeros)
	// * second other link rmb (can be all zeros)
	for i := range c.otherRMBs {
		c.otherRMBs[i].parse(buffer)
		buffer = buffer[13:]
	}

	// Rest of message is reserved
	c.res4 = buffer[0]
}

// String converts the confirm RKey message to a string
func (c *confirmRKey) String() string {
	var others string

	for i := range c.otherRMBs {
		others += fmt.Sprintf(", Other Link RMB %d: %s", i+1,
			&c.otherRMBs[i])
	}

	cFmt := "LLC Confirm RKey: Type: %d, Length: %d, Reply: %t, " +
		"Negative Response: %t, Configuration Retry: %t, " +
		"Number of Tokens: %d, This RKey: %d, This VAddr: %#x%s\n"
	return fmt.Sprintf(cFmt, c.Type, c.Length, c.reply, c.reject, c.retry,
		c.numTkns, c.rkey, c.vaddr, others)
}

// Reserved converts the confirm RKey message to a string including reserved
// fields
func (c *confirmRKey) Reserved() string {
	var others string

	for i := range c.otherRMBs {
		others += fmt.Sprintf("Other Link RMB %d: %s, ", i+1,
			&c.otherRMBs[i])
	}

	cFmt := "LLC Confirm RKey: Type: %d, Length: %d, Reserved: %#x, " +
		"Reply: %t, Reserved: %#x, Negative Response: %t, " +
		"Configuration Retry: %t, Reserved: %#x, " +
		"Number of Tokens: %d, This RKey: %d, This VAddr: %#x, " +
		"%sReserved: %#x\n"
	return fmt.Sprintf(cFmt, c.Type, c.Length, c.res1, c.reply, c.res2,
		c.reject, c.retry, c.res3, c.numTkns, c.rkey, c.vaddr, others,
		c.res4)
}

// parseConfirmRKey parses the LLC confirm RKey message in buffer
func parseConfirmRKey(buffer []byte) *confirmRKey {
	var confirm confirmRKey
	confirm.Parse(buffer)
	return &confirm
}
