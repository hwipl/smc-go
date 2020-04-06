package llc

import "fmt"

// ConfirmRKeyCont stores a LLC confirm rkey continuation message
type ConfirmRKeyCont struct {
	BaseMsg
	res1      byte
	Reply     bool
	res2      byte
	Reject    bool // negative response
	res3      byte
	NumTkns   uint8
	OtherRMBs [3]RMBSpec
	res4      byte
}

// Parse fills the confirmRKey fields from the confirm RKey continuation
// message in buffer
func (c *ConfirmRKeyCont) Parse(buffer []byte) {
	// TODO: merge with confirmRKey()?
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

	// Remainder of this byte is reserved
	c.res3 = buffer[0] & 0b00011111
	buffer = buffer[1:]

	// Number of tokens left is 1 byte
	c.NumTkns = buffer[0]
	buffer = buffer[1:]

	// other link rmb specifications are each 13 bytes
	// parse
	// * first other link rmb (can be all zeros)
	// * second other link rmb (can be all zeros)
	// * third other link rmb (can be all zeros)
	for i := range c.OtherRMBs {
		c.OtherRMBs[i].Parse(buffer)
		buffer = buffer[13:]
	}

	// Rest of message is reserved
	c.res4 = buffer[0]
}

// String converts the confirm RKey continuation message to a string
func (c *ConfirmRKeyCont) String() string {
	var others string

	for i := range c.OtherRMBs {
		others += fmt.Sprintf(", Other Link RMB %d: %s", i+1,
			&c.OtherRMBs[i])
	}

	cFmt := "LLC Confirm RKey Continuation: Type: %d, Length: %d, " +
		"Reply: %t, Negative Response: %t, Number of Tokens: %d%s\n"
	return fmt.Sprintf(cFmt, c.Type, c.Length, c.Reply, c.Reject, c.NumTkns,
		others)
}

// Reserved converts the confirm RKey continuation message to a string
// including reserved fields
func (c *ConfirmRKeyCont) Reserved() string {
	var others string

	for i := range c.OtherRMBs {
		others += fmt.Sprintf("Other Link RMB %d: %s, ", i+1,
			&c.OtherRMBs[i])
	}

	cFmt := "LLC Confirm RKey Continuation: Type: %d, Length: %d, " +
		"Reserved: %#x, Reply: %t, Reserved: %#x, " +
		"Negative Response: %t, Reserved: %#x, " +
		"Number of Tokens: %d, %sReserved: %#x\n"
	return fmt.Sprintf(cFmt, c.Type, c.Length, c.res1, c.Reply, c.res2,
		c.Reject, c.res3, c.NumTkns, others, c.res4)
}

// ParseConfirmRKeyCont parses the LLC confirm RKey Continuation message in
// buffer
func ParseConfirmRKeyCont(buffer []byte) *ConfirmRKeyCont {
	var confirmCont ConfirmRKeyCont
	confirmCont.Parse(buffer)
	return &confirmCont
}
