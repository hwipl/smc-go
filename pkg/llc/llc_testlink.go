package llc

import "fmt"

// testLink stores a LLC test link message
type testLink struct {
	BaseMsg
	res1     byte
	reply    bool
	res2     byte
	userData [16]byte
	res3     [24]byte
}

// Parse fills the testLink fields from the test link message in buffer
func (t *testLink) Parse(buffer []byte) {
	// init base message fields
	t.SetBaseMsg(buffer)
	buffer = buffer[2:]

	// Reserved 1 byte
	t.res1 = buffer[0]
	buffer = buffer[1:]

	// Reply is first bit in this byte
	t.reply = (buffer[0] & 0b10000000) > 0

	// Remainder of this byte is reserved
	t.res2 = buffer[0] & 0b01111111
	buffer = buffer[1:]

	// User data is 16 bytes
	copy(t.userData[:], buffer[0:16])
	buffer = buffer[16:]

	// Rest of message is reserved
	copy(t.res3[:], buffer[:])
}

// String converts the test link message to a string
func (t *testLink) String() string {
	tFmt := "LLC Test Link: Type %d, Length: %d, Reply: %t, " +
		"User Data: %#x\n"
	return fmt.Sprintf(tFmt, t.Type, t.Length, t.reply, t.userData)
}

// Reserved converts the test link message to a string including reserved
// fields
func (t *testLink) Reserved() string {
	tFmt := "LLC Test Link: Type %d, Length: %d, Reserved: %#x, " +
		"Reply: %t, Reserved: %#x, User Data: %#x, Reserved: %#x\n"
	return fmt.Sprintf(tFmt, t.Type, t.Length, t.res1, t.reply, t.res2,
		t.userData, t.res3)
}

// parseTestLink parses the LLC test link message in buffer
func parseTestLink(buffer []byte) *testLink {
	var test testLink
	test.Parse(buffer)
	return &test
}
