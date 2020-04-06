package llc

import "fmt"

// confirmRKeyCont stores a LLC confirm rkey continuation message
type confirmRKeyCont struct {
	baseMsg
	res1      byte
	reply     bool
	res2      byte
	reject    bool // negative response
	res3      byte
	numTkns   uint8
	otherRMBs [3]rmbSpec
	res4      byte
}

// Parse fills the confirmRKey fields from the confirm RKey continuation
// message in buffer
func (c *confirmRKeyCont) Parse(buffer []byte) {
	// TODO: merge with confirmRKey()?
	// init base message fields
	c.setBaseMsg(buffer)
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

	// Remainder of this byte is reserved
	c.res3 = buffer[0] & 0b00011111
	buffer = buffer[1:]

	// Number of tokens left is 1 byte
	c.numTkns = buffer[0]
	buffer = buffer[1:]

	// other link rmb specifications are each 13 bytes
	// parse
	// * first other link rmb (can be all zeros)
	// * second other link rmb (can be all zeros)
	// * third other link rmb (can be all zeros)
	for i := range c.otherRMBs {
		c.otherRMBs[i].parse(buffer)
		buffer = buffer[13:]
	}

	// Rest of message is reserved
	c.res4 = buffer[0]
}

// String converts the confirm RKey continuation message to a string
func (c *confirmRKeyCont) String() string {
	var others string

	for i := range c.otherRMBs {
		others += fmt.Sprintf(", Other Link RMB %d: %s", i+1,
			&c.otherRMBs[i])
	}

	cFmt := "LLC Confirm RKey Continuation: Type: %d, Length: %d, " +
		"Reply: %t, Negative Response: %t, Number of Tokens: %d%s\n"
	return fmt.Sprintf(cFmt, c.typ, c.length, c.reply, c.reject, c.numTkns,
		others)
}

// Reserved converts the confirm RKey continuation message to a string
// including reserved fields
func (c *confirmRKeyCont) Reserved() string {
	var others string

	for i := range c.otherRMBs {
		others += fmt.Sprintf("Other Link RMB %d: %s, ", i+1,
			&c.otherRMBs[i])
	}

	cFmt := "LLC Confirm RKey Continuation: Type: %d, Length: %d, " +
		"Reserved: %#x, Reply: %t, Reserved: %#x, " +
		"Negative Response: %t, Reserved: %#x, " +
		"Number of Tokens: %d, %sReserved: %#x\n"
	return fmt.Sprintf(cFmt, c.typ, c.length, c.res1, c.reply, c.res2,
		c.reject, c.res3, c.numTkns, others, c.res4)
}

// parseConfirmRKeyCont parses the LLC confirm RKey Continuation message in
// buffer
func parseConfirmRKeyCont(buffer []byte) *confirmRKeyCont {
	var confirmCont confirmRKeyCont
	confirmCont.Parse(buffer)
	return &confirmCont
}
