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

// opcode stores the bth opcode
type opcode uint8

// rcString converts the reliable connection opcode to a string
func (o opcode) rcString() string {
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
func (o opcode) ucString() string {
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
func (o opcode) rdString() string {
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
func (o opcode) udString() string {
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
func (o opcode) cnpString() string {
	op := int(o & 0b00011111)
	if op == 0b00000 {
		return "CNP"
	}
	return "Reserved"
}

// xrcString converts the extended reliable connection opcode to a string
func (o opcode) xrcString() string {
	// xrc strings are the same as rc strings
	return o.rcString()
}

// String converts the opcode to a string
func (o opcode) String() string {
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
	raw    []byte
	opcode opcode
	se     bool
	m      bool
	pad    uint8
	tver   uint8
	pkey   uint16
	fecn   bool
	becn   bool
	res1   byte
	destQP uint32
	a      bool
	res2   byte
	psn    uint32
}

// Parse fills the bth fields from the base transport header in buffer
func (b *BTH) Parse(buffer []byte) {
	// save raw message bytes, set internal type and length
	b.raw = make([]byte, len(buffer))
	copy(b.raw[:], buffer[:])

	// opcode is 1 byte
	b.opcode = opcode(buffer[0])
	buffer = buffer[1:]

	// solicited event is first bit in this byte
	b.se = (buffer[0] & 0b10000000) > 0

	// MigReq is the next bit in this byte
	b.m = (buffer[0] & 0b01000000) > 0

	// pad count is the next 2 bits in this byte
	b.pad = (buffer[0] & 0b00110000) >> 4

	// transport header version is last 4 bits in this byte
	b.tver = buffer[0] & 0b00001111
	buffer = buffer[1:]

	// partition key is 2 bytes
	b.pkey = binary.BigEndian.Uint16(buffer[0:2])
	buffer = buffer[2:]

	// FECN is first bit in this byte
	b.fecn = (buffer[0] & 0b10000000) > 0

	// BECN is next bit in this byte
	b.becn = (buffer[0] & 0b01000000) > 0

	// Reserved are the last 6 bits in this byte
	b.res1 = buffer[0] & 0b00111111
	buffer = buffer[1:]

	// destination QP number is 3 bytes
	b.destQP = uint32(buffer[0]) << 16
	b.destQP |= uint32(buffer[1]) << 8
	b.destQP |= uint32(buffer[2])
	buffer = buffer[3:]

	// AckReq is first bit in this byte
	b.a = (buffer[0] & 0b10000000) > 0

	// Reserved are the last 7 bits in this byte
	b.res2 = buffer[0] & 0b01111111
	buffer = buffer[1:]

	// Packet Sequence Number is 3 bytes
	b.psn = uint32(buffer[0]) << 16
	b.psn |= uint32(buffer[1]) << 8
	b.psn |= uint32(buffer[2])
}

// String converts the base transport header to a string
func (b *BTH) String() string {
	bfmt := "BTH: OpCode: %s, SE: %t, M: %t, Pad: %d, TVer: %d, " +
		"PKey: %d, FECN: %t, BECN: %t, DestQP: %d, A: %t, PSN: %d\n"
	return fmt.Sprintf(bfmt, b.opcode, b.se, b.m, b.pad, b.tver, b.pkey,
		b.fecn, b.becn, b.destQP, b.a, b.psn)
}

// Reserved converts the base transport header to a string
func (b *BTH) Reserved() string {
	bfmt := "BTH: OpCode: %s, SE: %t, M: %t, Pad: %d, TVer: %d, " +
		"PKey: %d, FECN: %t, BECN: %t, Res: %#x, DestQP: %d, " +
		"A: %t, Res: %#x, PSN: %d\n"
	return fmt.Sprintf(bfmt, b.opcode, b.se, b.m, b.pad, b.tver, b.pkey,
		b.fecn, b.becn, b.res1, b.destQP, b.a, b.res2, b.psn)
}

// Hex converts the message to a hex dump string
func (b *BTH) Hex() string {
	return hex.Dump(b.raw)
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
