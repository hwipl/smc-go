package llc

// Message is an interface for messages
type Message interface {
	Parse(buffer []byte)
	String() string
	Reserved() string
	Hex() string
	GetType() int
}
