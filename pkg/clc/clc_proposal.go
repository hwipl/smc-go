package clc

import (
	"encoding/binary"
	"fmt"
	"log"
	"net"
)

const (
	ProposalLen      = 52 // minimum length
	IPv6PrefixLen    = 17
	SMCDIPAreaOffset = 40
)

// IPv6Prefix stores a SMC IPv6 Prefix
type IPv6Prefix struct {
	prefix    net.IP
	prefixLen uint8
}

// String converts ipv6Prefix to a string
func (p IPv6Prefix) String() string {
	return fmt.Sprintf("%s/%d", p.prefix, p.prefixLen)
}

// Proposal stores a CLC Proposal message
type Proposal struct {
	raw
	Header
	SenderPeerID peerID           // unique system id
	IBGID        net.IP           // gid of ib_device port
	IBMAC        net.HardwareAddr // mac of ib_device port
	IPAreaOffset uint16           // offset to IP address info area

	// Optional SMC-D info
	SMCDGID  uint64 // ISM GID of requestor
	reserved [32]byte

	// IP/Prefix info
	Prefix          net.IP // subnet mask (rather prefix)
	PrefixLen       uint8  // number of significant bits in mask
	reserved2       [2]byte
	IPv6PrefixesCnt uint8 // number of IPv6 prefixes in prefix array
	IPv6Prefixes    []IPv6Prefix

	trailer
}

// String converts the CLC Proposal message to a string
func (p *Proposal) String() string {
	if p == nil {
		return "n/a"
	}

	// smc-d info
	smcdInfo := ""
	if p.IPAreaOffset == SMCDIPAreaOffset {
		smcdInfo = fmt.Sprintf("SMC-D GID: %d, ", p.SMCDGID)
	}

	// ipv6 prefixes
	ipv6Prefixes := ""
	for _, prefix := range p.IPv6Prefixes {
		ipv6Prefixes += fmt.Sprintf(", IPv6 Prefix: %s", prefix)
	}

	proposalFmt := "%s, Peer ID: %s, SMC-R GID: %s, RoCE MAC: %s, " +
		"IP Area Offset: %d, %sIPv4 Prefix: %s/%d, " +
		"IPv6 Prefix Count: %d%s, Trailer: %s"
	return fmt.Sprintf(proposalFmt, p.Header.String(), p.SenderPeerID,
		p.IBGID, p.IBMAC, p.IPAreaOffset, smcdInfo, p.Prefix,
		p.PrefixLen, p.IPv6PrefixesCnt, ipv6Prefixes, p.trailer)
}

// Reserved converts the CLC Proposal message to a string including reserved
// message fields
func (p *Proposal) Reserved() string {
	if p == nil {
		return "n/a"
	}

	// smc-d info
	smcdInfo := ""
	if p.IPAreaOffset == SMCDIPAreaOffset {
		smcdInfo = fmt.Sprintf("SMC-D GID: %d, Reserved: %#x, ",
			p.SMCDGID, p.reserved)
	}
	// ipv6 prefixes
	ipv6Prefixes := ""
	for _, prefix := range p.IPv6Prefixes {
		ipv6Prefixes += fmt.Sprintf(", IPv6 Prefix: %s", prefix)
	}

	proposalFmt := "%s, Peer ID: %s, SMC-R GID: %s, RoCE MAC: %s, " +
		"IP Area Offset: %d, %sIPv4 Prefix: %s/%d, Reserved: %#x, " +
		"IPv6 Prefix Count: %d%s, Trailer: %s"
	return fmt.Sprintf(proposalFmt, p.Header.Reserved(), p.SenderPeerID,
		p.IBGID, p.IBMAC, p.IPAreaOffset, smcdInfo, p.Prefix,
		p.PrefixLen, p.reserved2, p.IPv6PrefixesCnt, ipv6Prefixes,
		p.trailer)
}

// Parse parses the CLC Proposal message in buf
func (p *Proposal) Parse(buf []byte) {
	// save raw message bytes
	p.raw.Parse(buf)

	// parse CLC header
	p.Header.Parse(buf)

	// check if message is long enough
	if p.Length < ProposalLen {
		log.Println("Error parsing CLC Proposal: message too short")
		errDump(buf[:p.Length])
		return
	}

	// skip clc header
	skip := HeaderLen

	// sender peer ID
	copy(p.SenderPeerID[:], buf[skip:skip+peerIDLen])
	skip += peerIDLen

	// ib GID is an IPv6 address
	p.IBGID = make(net.IP, net.IPv6len)
	copy(p.IBGID[:], buf[skip:skip+net.IPv6len])
	skip += net.IPv6len

	// ib MAC is a 6 byte MAC address
	p.IBMAC = make(net.HardwareAddr, 6)
	copy(p.IBMAC[:], buf[skip:skip+6])
	skip += 6

	// offset to ip area
	p.IPAreaOffset = binary.BigEndian.Uint16(buf[skip : skip+2])
	skip += 2

	// Optional SMC-D info
	if p.IPAreaOffset == SMCDIPAreaOffset {
		// smcd GID
		p.SMCDGID = binary.BigEndian.Uint64(buf[skip : skip+8])
		skip += 8

		// reserved
		copy(p.reserved[:], buf[skip:skip+32])
		skip += 32
	} else {
		skip += int(p.IPAreaOffset)
	}

	// make sure we do not read outside the message
	if int(p.Length)-skip < net.IPv4len+1+2+1+trailerLen {
		log.Println("Error parsing CLC Proposal: " +
			"IP Area Offset too big")
		errDump(buf[:p.Length])
		return
	}

	// IP/prefix is an IPv4 address
	p.Prefix = make(net.IP, net.IPv4len)
	copy(p.Prefix[:], buf[skip:skip+net.IPv4len])
	skip += net.IPv4len

	// prefix length
	p.PrefixLen = uint8(buf[skip])
	skip++

	// reserved
	copy(p.reserved2[:], buf[skip:skip+2])
	skip += 2

	// ipv6 prefix count
	p.IPv6PrefixesCnt = uint8(buf[skip])

	// parse ipv6 prefixes
	for i := uint8(0); i < p.IPv6PrefixesCnt; i++ {
		// skip prefix count or last prefix length
		skip++

		// make sure we are still inside the clc message
		if int(p.Length)-skip < IPv6PrefixLen+trailerLen {
			log.Println("Error parsing CLC Proposal: " +
				"IPv6 prefix count too big")
			errDump(buf[:p.Length])
			break
		}
		// create new ipv6 prefix entry
		ip6prefix := IPv6Prefix{}

		// parse prefix and fill prefix entry
		ip6prefix.prefix = make(net.IP, net.IPv6len)
		copy(ip6prefix.prefix[:], buf[skip:skip+net.IPv6len])
		skip += net.IPv6len

		// parse prefix length and fill prefix entry
		ip6prefix.prefixLen = uint8(buf[skip])

		// add to ipv6 prefixes
		p.IPv6Prefixes = append(p.IPv6Prefixes, ip6prefix)
	}

	// save trailer
	p.trailer.Parse(buf)
}
