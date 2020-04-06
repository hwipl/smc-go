package llc

const (
	// TypeOther is the internal message type for other/non-LLC messages
	TypeOther = 0x101
)

// Other stores an Other message
type Other struct {
	BaseMsg
}

// Parse fills the other fields from the other message in buffer
func (o *Other) Parse(buffer []byte) {
	o.setRaw(buffer)
	o.Type = TypeOther
	o.Length = len(buffer)
}

// String converts the other message into a string
func (o *Other) String() string {
	return "Other Payload\n"
}

// Reserved converts the other message into a string including reserved fields
func (o *Other) Reserved() string {
	return o.String()
}

// ParseOther parses the other message in buffer
func ParseOther(buffer []byte) *Other {
	var o Other
	o.Parse(buffer)
	return &o
}
