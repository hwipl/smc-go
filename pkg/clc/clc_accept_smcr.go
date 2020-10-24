package clc

import (
	"encoding/binary"
	"fmt"
	"log"
	"net"
)

const (
	// AcceptSMCRLen is the length of a CLC SMC-R Accept message
	AcceptSMCRLen = 68
)

// QPMTU stores a SMC QP MTU
type QPMTU uint8

// String converts qpMTU to a string
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

// AcceptSMCR stores a CLC SMC-R Accept message
type AcceptSMCR struct {
	Raw
	Header
	SenderPeerID   PeerID           // unique system id
	IBGID          net.IP           // gid of ib_device port
	IBMAC          net.HardwareAddr // mac of ib_device port
	QPN            int              // QP number
	RMBRKey        uint32           // RMB rkey
	RMBEIdx        uint8            // Index of RMBE in RMB
	RMBEAlertToken uint32           // unique connection id
	RMBESize       RMBESize         // 4 bits buf size (compressed)
	QPMTU          QPMTU            // 4 bits QP mtu
	reserved       byte
	RMBDMAAddr     uint64 // RMB virtual address
	reserved2      byte
	PSN            int // packet sequence number
	Trailer
}

// String converts the CLC SMC-R Accept message to a string
func (ac *AcceptSMCR) String() string {
	if ac == nil {
		return "n/a"
	}

	acFmt := "%s, Peer ID: %s, SMC-R GID: %s, RoCE MAC: %s, " +
		"QP Number: %d, RMB RKey: %d, RMBE Index: %d, " +
		"RMBE Alert Token: %d, RMBE Size: %s, QP MTU: %s, " +
		"RMB Virtual Address: %#x, Packet Sequence Number: %d, " +
		"Trailer: %s"
	return fmt.Sprintf(acFmt, ac.Header.String(), ac.SenderPeerID,
		ac.IBGID, ac.IBMAC, ac.QPN, ac.RMBRKey, ac.RMBEIdx,
		ac.RMBEAlertToken, ac.RMBESize, ac.QPMTU, ac.RMBDMAAddr,
		ac.PSN, ac.Trailer)
}

// Reserved converts the CLC SMC-R Accept message to a string including
// reserved message fields
func (ac *AcceptSMCR) Reserved() string {
	if ac == nil {
		return "n/a"
	}

	acFmt := "%s, Peer ID: %s, SMC-R GID: %s, RoCE MAC: %s, " +
		"QP Number: %d, RMB RKey: %d, RMBE Index: %d, " +
		"RMBE Alert Token: %d, RMBE Size: %s, QP MTU: %s, " +
		"Reserved: %#x, RMB Virtual Address: %#x, " +
		"Reserved: %#x, Packet Sequence Number: %d, Trailer: %s"
	return fmt.Sprintf(acFmt, ac.Header.Reserved(), ac.SenderPeerID,
		ac.IBGID, ac.IBMAC, ac.QPN, ac.RMBRKey, ac.RMBEIdx,
		ac.RMBEAlertToken, ac.RMBESize, ac.QPMTU, ac.reserved,
		ac.RMBDMAAddr, ac.reserved2, ac.PSN, ac.Trailer)
}

// Parse parses the SMC-R Accept message in buf
func (ac *AcceptSMCR) Parse(buf []byte) {
	// save raw message bytes
	ac.Raw.Parse(buf)

	// parse CLC header
	ac.Header.Parse(buf)

	// check if message is long enough
	if ac.Length < AcceptSMCRLen {
		err := "Error parsing CLC Accept: message too short"
		if ac.Type == TypeConfirm {
			err = "Error parsing CLC Confirm: message too short"
		}
		log.Println(err)
		errDump(buf[:ac.Length])
		return
	}

	// skip clc header
	buf = buf[HeaderLen:]

	// sender peer ID
	copy(ac.SenderPeerID[:], buf[:PeerIDLen])
	buf = buf[PeerIDLen:]

	// ib GID is an IPv6 Address
	ac.IBGID = make(net.IP, net.IPv6len)
	copy(ac.IBGID[:], buf[:net.IPv6len])
	buf = buf[net.IPv6len:]

	// ib MAC is a 6 byte MAC address
	ac.IBMAC = make(net.HardwareAddr, 6)
	copy(ac.IBMAC[:], buf[:6])
	buf = buf[6:]

	// QP number is 3 bytes
	ac.QPN = int(buf[0]) << 16
	ac.QPN |= int(buf[1]) << 8
	ac.QPN |= int(buf[2])
	buf = buf[3:]

	// rmb Rkey
	ac.RMBRKey = binary.BigEndian.Uint32(buf[:4])
	buf = buf[4:]

	// rmbe Idx
	ac.RMBEIdx = uint8(buf[0])
	buf = buf[1:]

	// rmbe alert token
	ac.RMBEAlertToken = binary.BigEndian.Uint32(buf[:4])
	buf = buf[4:]

	// 1 byte bitfield: rmbe size (4 bits) and qp mtu (4 bits)
	ac.RMBESize = RMBESize((uint8(buf[0]) & 0b11110000) >> 4)
	ac.QPMTU = QPMTU(uint8(buf[0]) & 0b00001111)
	buf = buf[1:]

	// reserved
	ac.reserved = buf[0]
	buf = buf[1:]

	// rmb DMA addr
	ac.RMBDMAAddr = binary.BigEndian.Uint64(buf[:8])
	buf = buf[8:]

	// reserved
	ac.reserved2 = buf[0]
	buf = buf[1:]

	// Packet Sequence Number is 3 bytes
	ac.PSN = int(buf[0]) << 16
	ac.PSN |= int(buf[1]) << 8
	ac.PSN |= int(buf[2])
	buf = buf[3:]

	// save trailer
	ac.Trailer.Parse(buf)
}
