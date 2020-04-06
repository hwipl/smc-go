package llc

const (
	// TypeOther is the internal message type for other/non-LLC messages
	TypeOther = 0x101
)

// other stores an other message
type other struct {
	baseMsg
}

// Parse fills the other fields from the other message in buffer
func (o *other) Parse(buffer []byte) {
	o.setRaw(buffer)
	o.typ = TypeOther
	o.length = len(buffer)
}

// String converts the other message into a string
func (o *other) String() string {
	return "Other Payload\n"
}

// Reserved converts the other message into a string including reserved fields
func (o *other) Reserved() string {
	return o.String()
}

// parseOther parses the other message in buffer
func parseOther(buffer []byte) *other {
	var o other
	o.Parse(buffer)
	return &o
}
