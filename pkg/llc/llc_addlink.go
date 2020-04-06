package llc

import (
	"fmt"
	"net"
)

// QPMTU stores the compressed MTU of a QP, taken from smc-clc
type QPMTU uint8

// String converts qpMTU to a string including the uncompressed MTU, taken from
// smc-clc
func (m QPMTU) String() string {
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

// AddLinkRsnCode stores the reason code of LLC add link messages
type AddLinkRsnCode uint8

// String converts the reason code to a string
func (r AddLinkRsnCode) String() string {
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

// AddLink stores a LLC add link message
type AddLink struct {
	BaseMsg
	res1      byte
	RsnCode   AddLinkRsnCode
	Reply     bool
	Reject    bool
	res2      byte
	SenderMAC net.HardwareAddr
	res5      [2]byte // not in the RFC
	SenderGID net.IP
	SenderQP  uint32
	Link      uint8
	res3      byte
	MTU       QPMTU
	PSN       uint32
	res4      [10]byte
}

// Parse fills the addLink fields from the LLC add link message in buffer
func (a *AddLink) Parse(buffer []byte) {
	// init base message fields
	a.SetBaseMsg(buffer)
	buffer = buffer[2:]

	// Reserved are first 4 bits in this byte
	a.res1 = buffer[0] >> 4

	// Reason Code are the last 4 bits in this byte
	a.RsnCode = AddLinkRsnCode(buffer[0] & 0b00001111)
	buffer = buffer[1:]

	// Reply flag is the first bit in this byte
	a.Reply = (buffer[0] & 0b10000000) > 0

	// Rejection flag is the next bit in this byte
	a.Reject = (buffer[0] & 0b01000000) > 0

	// Reserved are the last 6 bits in this byte
	a.res2 = buffer[0] & 0b00111111
	buffer = buffer[1:]

	// sender MAC is a 6 byte MAC address
	a.SenderMAC = make(net.HardwareAddr, 6)
	copy(a.SenderMAC[:], buffer[0:6])
	buffer = buffer[6:]

	// in the linux code, there are 2 more reserved bytes here that are not
	// in the RFC
	copy(a.res5[:], buffer[0:2])
	buffer = buffer[2:]

	// sender GID is an 16 bytes IPv6 address
	a.SenderGID = make(net.IP, net.IPv6len)
	copy(a.SenderGID[:], buffer[0:16])
	buffer = buffer[16:]

	// QP number is 3 bytes
	a.SenderQP = uint32(buffer[0]) << 16
	a.SenderQP |= uint32(buffer[1]) << 8
	a.SenderQP |= uint32(buffer[2])
	buffer = buffer[3:]

	// Link is 1 byte
	a.Link = buffer[0]
	buffer = buffer[1:]

	// Reserved are the first 4 bits in this byte
	a.res3 = buffer[0] >> 4

	// MTU are the last 4 bits in this byte
	a.MTU = QPMTU(buffer[0] & 0b00001111)
	buffer = buffer[1:]

	// initial Packet Sequence Number is 3 bytes
	a.PSN = uint32(buffer[0]) << 16
	a.PSN |= uint32(buffer[1]) << 8
	a.PSN |= uint32(buffer[2])
	buffer = buffer[3:]

	// Rest of message is reserved
	copy(a.res4[:], buffer[:])
}

// String converts the LLC add link message to string
func (a *AddLink) String() string {
	aFmt := "LLC Add Link: Type: %d, Length: %d, Reason Code: %s, " +
		"Reply: %t, Rejection: %t, Sender MAC: %s, Sender GID: %s, " +
		"Sender QP: %d, Link: %d, MTU: %s, Initial PSN: %d\n"
	return fmt.Sprintf(aFmt, a.Type, a.Length, a.RsnCode, a.Reply, a.Reject,
		a.SenderMAC, a.SenderGID, a.SenderQP, a.Link, a.MTU, a.PSN)
}

// Reserved converts the LLC add link message to string including reserved
// fields
func (a *AddLink) Reserved() string {
	aFmt := "LLC Add Link: Type: %d, Length: %d, Reserved: %#x, " +
		"Reason Code: %s, Reply: %t, Rejection: %t, Reserved: %#x, " +
		"Sender MAC: %s, Sender GID: %s, Sender QP: %d, Link: %d, " +
		"Reserved: %#x, MTU: %s, Initial PSN: %d, Reserved: %#x\n"
	return fmt.Sprintf(aFmt, a.Type, a.Length, a.res1, a.RsnCode, a.Reply,
		a.Reject, a.res2, a.SenderMAC, a.SenderGID, a.SenderQP, a.Link,
		a.res3, a.MTU, a.PSN, a.res4)
}

// ParseAddLink parses and prints the LLC add link message in buffer
func ParseAddLink(buffer []byte) *AddLink {
	var add AddLink
	add.Parse(buffer)
	return &add
}
