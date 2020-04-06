package llc

import "fmt"

// TestLink stores a LLC test link message
type TestLink struct {
	BaseMsg
	res1     byte
	Reply    bool
	res2     byte
	UserData [16]byte
	res3     [24]byte
}

// Parse fills the testLink fields from the test link message in buffer
func (t *TestLink) Parse(buffer []byte) {
	// init base message fields
	t.SetBaseMsg(buffer)
	buffer = buffer[2:]

	// Reserved 1 byte
	t.res1 = buffer[0]
	buffer = buffer[1:]

	// Reply is first bit in this byte
	t.Reply = (buffer[0] & 0b10000000) > 0

	// Remainder of this byte is reserved
	t.res2 = buffer[0] & 0b01111111
	buffer = buffer[1:]

	// User data is 16 bytes
	copy(t.UserData[:], buffer[0:16])
	buffer = buffer[16:]

	// Rest of message is reserved
	copy(t.res3[:], buffer[:])
}

// String converts the test link message to a string
func (t *TestLink) String() string {
	tFmt := "LLC Test Link: Type %d, Length: %d, Reply: %t, " +
		"User Data: %#x\n"
	return fmt.Sprintf(tFmt, t.Type, t.Length, t.Reply, t.UserData)
}

// Reserved converts the test link message to a string including reserved
// fields
func (t *TestLink) Reserved() string {
	tFmt := "LLC Test Link: Type %d, Length: %d, Reserved: %#x, " +
		"Reply: %t, Reserved: %#x, User Data: %#x, Reserved: %#x\n"
	return fmt.Sprintf(tFmt, t.Type, t.Length, t.res1, t.Reply, t.res2,
		t.UserData, t.res3)
}

// ParseTestLink parses the LLC test link message in buffer
func ParseTestLink(buffer []byte) *TestLink {
	var test TestLink
	test.Parse(buffer)
	return &test
}
