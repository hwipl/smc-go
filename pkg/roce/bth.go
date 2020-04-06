package roce

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
)

const (
	// internal message type
	typeBTH = 0x103

	// BTHNextHeader is the next header value in the GRH
	BTHNextHeader = 0x1B

	// BTHLen is the length of the bth
	BTHLen = 12
)

// Opcode stores the bth Opcode
type Opcode uint8

// rcString converts the reliable connection opcode to a string
func (o Opcode) rcString() string {
	var rcStrings = [...]string{
		"SEND First",
		"SEND Middle",
		"SEND Last",
		"SEND Last with Immediate",
		"SEND Only",
		"SEND Only with Immediate",
		"RDMA WRITE First",
		"RDMA WRITE Middle",
		"RDMA WRITE Last",
		"RDMA WRITE Last with Immediate",
		"RDMA WRITE Only",
		"RDMA WRITE Only with Immediate",
		"RDMA READ Request",
		"RDMA READ response First",
		"RDMA READ response Middle",
		"RDMA READ response Last",
		"RDMA READ response Only",
		"Acknowledge",
		"ATOMIC Acknowledge",
		"CmpSwap",
		"FetchAdd",
		"Reserved",
		"SEND Last with Invalidate",
		"SEND Only with Invalidate",
	}
	// lookup last 5 bits of opcode in rcStrings to get the string
	op := int(o & 0b00011111)
	if op < len(rcStrings) {
		return rcStrings[op]
	}
	return "Reserved"
}

// ucString converts the unreliable connection opcode to a string
func (o Opcode) ucString() string {
	var ucStrings = [...]string{
		"SEND First",
		"SEND Middle",
		"SEND Last",
		"SEND Last with Immediate",
		"SEND Only",
		"SEND Only with Immediate",
		"RDMA WRITE First",
		"RDMA WRITE Middle",
		"RDMA WRITE Last",
		"RDMA WRITE Last with Immediate",
		"RDMA WRITE Only",
		"RDMA WRITE Only with Immediate",
	}
	// lookup last 5 bits of opcode in ucStrings to get the string
	op := int(o & 0b00011111)
	if op < len(ucStrings) {
		return ucStrings[op]
	}
	return "Reserved"
}

// rdString converts the reliable datagram opcode to a string
func (o Opcode) rdString() string {
	var rdStrings = [...]string{
		"SEND First",
		"SEND Middle",
		"SEND Last",
		"SEND Last with Immediate",
		"SEND Only",
		"SEND Only with Immediate",
		"RDMA WRITE First",
		"RDMA WRITE Middle",
		"RDMA WRITE Last",
		"RDMA WRITE Last with Immediate",
		"RDMA WRITE Only",
		"RDMA WRITE Only with Immediate",
		"RDMA READ Request",
		"RDMA READ response First",
		"RDMA READ response Middle",
		"RDMA READ response Last",
		"RDMA READ response Only",
		"Acknowledge",
		"ATOMIC Acknowledge",
		"CmpSwap",
		"FetchAdd",
		"RESYNC",
	}
	// lookup last 5 bits of opcode in rdStrings to get the string
	op := int(o & 0b00011111)
	if op < len(rdStrings) {
		return rdStrings[op]
	}
	return "Reserved"
}

// udString converts the unreliable datagram opcode to a string
func (o Opcode) udString() string {
	var udStrings = [...]string{
		"Reserved",
		"Reserved",
		"Reserved",
		"Reserved",
		"SEND Only",
		"SEND Only with Immediate",
	}
	// lookup last 5 bits of opcode in udStrings to get the string
	op := int(o & 0b00011111)
	if op < len(udStrings) {
		return udStrings[op]
	}
	return "Reserved"
}

// cnpString converts the CNP opcode to a string
func (o Opcode) cnpString() string {
	op := int(o & 0b00011111)
	if op == 0b00000 {
		return "CNP"
	}
	return "Reserved"
}

// xrcString converts the extended reliable connection opcode to a string
func (o Opcode) xrcString() string {
	// xrc strings are the same as rc strings
	return o.rcString()
}

// String converts the opcode to a string
func (o Opcode) String() string {
	var op string
	switch o >> 5 {
	case 0b000:
		// Reliable Connection (RC)
		op = "RC " + o.rcString()
	case 0b001:
		// Unreliable Connection (UC)
		op = "UC " + o.ucString()
	case 0b010:
		// Reliable Datagram (RD)
		op = "RD " + o.rdString()
	case 0b011:
		// Unreliable Datagram (UD)
		op = "UD " + o.udString()
	case 0b100:
		// CNP
		op = "CNP " + o.cnpString()
	case 0b101:
		// Extended Reliable Connection (XRC)
		op = "XRC " + o.xrcString()
	default:
		// Manufacturer Specific OpCodes
		op = "Manufacturer Specific"
	}
	return fmt.Sprintf("%#b (%s)", o, op)
}

// BTH stores an ib base transport header
type BTH struct {
	Raw    []byte
	Opcode Opcode
	SE     bool
	M      bool
	Pad    uint8
	TVer   uint8
	PKey   uint16
	FECN   bool
	BECN   bool
	res1   byte
	DestQP uint32
	A      bool
	res2   byte
	PSN    uint32
}

// Parse fills the bth fields from the base transport header in buffer
func (b *BTH) Parse(buffer []byte) {
	// save raw message bytes, set internal type and length
	b.Raw = make([]byte, len(buffer))
	copy(b.Raw[:], buffer[:])

	// opcode is 1 byte
	b.Opcode = Opcode(buffer[0])
	buffer = buffer[1:]

	// solicited event is first bit in this byte
	b.SE = (buffer[0] & 0b10000000) > 0

	// MigReq is the next bit in this byte
	b.M = (buffer[0] & 0b01000000) > 0

	// pad count is the next 2 bits in this byte
	b.Pad = (buffer[0] & 0b00110000) >> 4

	// transport header version is last 4 bits in this byte
	b.TVer = buffer[0] & 0b00001111
	buffer = buffer[1:]

	// partition key is 2 bytes
	b.PKey = binary.BigEndian.Uint16(buffer[0:2])
	buffer = buffer[2:]

	// FECN is first bit in this byte
	b.FECN = (buffer[0] & 0b10000000) > 0

	// BECN is next bit in this byte
	b.BECN = (buffer[0] & 0b01000000) > 0

	// Reserved are the last 6 bits in this byte
	b.res1 = buffer[0] & 0b00111111
	buffer = buffer[1:]

	// destination QP number is 3 bytes
	b.DestQP = uint32(buffer[0]) << 16
	b.DestQP |= uint32(buffer[1]) << 8
	b.DestQP |= uint32(buffer[2])
	buffer = buffer[3:]

	// AckReq is first bit in this byte
	b.A = (buffer[0] & 0b10000000) > 0

	// Reserved are the last 7 bits in this byte
	b.res2 = buffer[0] & 0b01111111
	buffer = buffer[1:]

	// Packet Sequence Number is 3 bytes
	b.PSN = uint32(buffer[0]) << 16
	b.PSN |= uint32(buffer[1]) << 8
	b.PSN |= uint32(buffer[2])
}

// String converts the base transport header to a string
func (b *BTH) String() string {
	bfmt := "BTH: OpCode: %s, SE: %t, M: %t, Pad: %d, TVer: %d, " +
		"PKey: %d, FECN: %t, BECN: %t, DestQP: %d, A: %t, PSN: %d\n"
	return fmt.Sprintf(bfmt, b.Opcode, b.SE, b.M, b.Pad, b.TVer, b.PKey,
		b.FECN, b.BECN, b.DestQP, b.A, b.PSN)
}

// Reserved converts the base transport header to a string
func (b *BTH) Reserved() string {
	bfmt := "BTH: OpCode: %s, SE: %t, M: %t, Pad: %d, TVer: %d, " +
		"PKey: %d, FECN: %t, BECN: %t, Res: %#x, DestQP: %d, " +
		"A: %t, Res: %#x, PSN: %d\n"
	return fmt.Sprintf(bfmt, b.Opcode, b.SE, b.M, b.Pad, b.TVer, b.PKey,
		b.FECN, b.BECN, b.res1, b.DestQP, b.A, b.res2, b.PSN)
}

// Hex converts the message to a hex dump string
func (b *BTH) Hex() string {
	return hex.Dump(b.Raw)
}

// GetType returns the type of the message
func (b *BTH) GetType() int {
	return typeBTH
}

// ParseBTH parses the BTH header in buffer
func ParseBTH(buffer []byte) *BTH {
	var b BTH
	b.Parse(buffer)
	return &b
}
