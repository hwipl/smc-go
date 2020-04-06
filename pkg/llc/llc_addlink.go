package llc

import (
	"fmt"
	"net"
)

// qpMTU stores the compressed MTU of a QP, taken from smc-clc
type qpMTU uint8

// String converts qpMTU to a string including the uncompressed MTU, taken from
// smc-clc
func (m qpMTU) String() string {
	var mtu string

	switch m {
	case 1:
		mtu = "256"
	case 2:
		mtu = "512"
	case 3:
		mtu = "1024"
	case 4:
		mtu = "2048"
	case 5:
		mtu = "4096"
	default:
		mtu = "reserved"
	}

	return fmt.Sprintf("%d (%s)", m, mtu)
}

// addLinkRsnCode stores the reason code of LLC add link messages
type addLinkRsnCode uint8

// String converts the reason code to a string
func (r addLinkRsnCode) String() string {
	var rsn string

	switch r {
	case 1:
		rsn = "no alternate path available"
	case 2:
		rsn = "invalid MTU value specified"
	default:
		rsn = "unknown"
	}

	return fmt.Sprintf("%d (%s)", r, rsn)
}

// addLink stores a LLC add link message
type addLink struct {
	baseMsg
	res1      byte
	rsnCode   addLinkRsnCode
	reply     bool
	reject    bool
	res2      byte
	senderMAC net.HardwareAddr
	res5      [2]byte // not in the RFC
	senderGID net.IP
	senderQP  uint32
	link      uint8
	res3      byte
	mtu       qpMTU
	psn       uint32
	res4      [10]byte
}

// Parse fills the addLink fields from the LLC add link message in buffer
func (a *addLink) Parse(buffer []byte) {
	// init base message fields
	a.setBaseMsg(buffer)
	buffer = buffer[2:]

	// Reserved are first 4 bits in this byte
	a.res1 = buffer[0] >> 4

	// Reason Code are the last 4 bits in this byte
	a.rsnCode = addLinkRsnCode(buffer[0] & 0b00001111)
	buffer = buffer[1:]

	// Reply flag is the first bit in this byte
	a.reply = (buffer[0] & 0b10000000) > 0

	// Rejection flag is the next bit in this byte
	a.reject = (buffer[0] & 0b01000000) > 0

	// Reserved are the last 6 bits in this byte
	a.res2 = buffer[0] & 0b00111111
	buffer = buffer[1:]

	// sender MAC is a 6 byte MAC address
	a.senderMAC = make(net.HardwareAddr, 6)
	copy(a.senderMAC[:], buffer[0:6])
	buffer = buffer[6:]

	// in the linux code, there are 2 more reserved bytes here that are not
	// in the RFC
	copy(a.res5[:], buffer[0:2])
	buffer = buffer[2:]

	// sender GID is an 16 bytes IPv6 address
	a.senderGID = make(net.IP, net.IPv6len)
	copy(a.senderGID[:], buffer[0:16])
	buffer = buffer[16:]

	// QP number is 3 bytes
	a.senderQP = uint32(buffer[0]) << 16
	a.senderQP |= uint32(buffer[1]) << 8
	a.senderQP |= uint32(buffer[2])
	buffer = buffer[3:]

	// Link is 1 byte
	a.link = buffer[0]
	buffer = buffer[1:]

	// Reserved are the first 4 bits in this byte
	a.res3 = buffer[0] >> 4

	// MTU are the last 4 bits in this byte
	a.mtu = qpMTU(buffer[0] & 0b00001111)
	buffer = buffer[1:]

	// initial Packet Sequence Number is 3 bytes
	a.psn = uint32(buffer[0]) << 16
	a.psn |= uint32(buffer[1]) << 8
	a.psn |= uint32(buffer[2])
	buffer = buffer[3:]

	// Rest of message is reserved
	copy(a.res4[:], buffer[:])
}

// String converts the LLC add link message to string
func (a *addLink) String() string {
	aFmt := "LLC Add Link: Type: %d, Length: %d, Reason Code: %s, " +
		"Reply: %t, Rejection: %t, Sender MAC: %s, Sender GID: %s, " +
		"Sender QP: %d, Link: %d, MTU: %s, Initial PSN: %d\n"
	return fmt.Sprintf(aFmt, a.typ, a.length, a.rsnCode, a.reply, a.reject,
		a.senderMAC, a.senderGID, a.senderQP, a.link, a.mtu, a.psn)
}

// Reserved converts the LLC add link message to string including reserved
// fields
func (a *addLink) Reserved() string {
	aFmt := "LLC Add Link: Type: %d, Length: %d, Reserved: %#x, " +
		"Reason Code: %s, Reply: %t, Rejection: %t, Reserved: %#x, " +
		"Sender MAC: %s, Sender GID: %s, Sender QP: %d, Link: %d, " +
		"Reserved: %#x, MTU: %s, Initial PSN: %d, Reserved: %#x\n"
	return fmt.Sprintf(aFmt, a.typ, a.length, a.res1, a.rsnCode, a.reply,
		a.reject, a.res2, a.senderMAC, a.senderGID, a.senderQP, a.link,
		a.res3, a.mtu, a.psn, a.res4)
}

// parseAddLink parses and prints the LLC add link message in buffer
func parseAddLink(buffer []byte) *addLink {
	var add addLink
	add.Parse(buffer)
	return &add
}
