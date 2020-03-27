package clc

import (
	"encoding/binary"
	"fmt"
	"log"
	"net"
)

const (
	proposalLen      = 52 // minimum length
	ipv6PrefixLen    = 17
	smcdIPAreaOffset = 40
)

// ipv6Prefix stores a SMC IPv6 Prefix
type ipv6Prefix struct {
	prefix    net.IP
	prefixLen uint8
}

// String converts ipv6Prefix to a string
func (p ipv6Prefix) String() string {
	return fmt.Sprintf("%s/%d", p.prefix, p.prefixLen)
}

// proposal stores a CLC Proposal message
type proposal struct {
	raw
	header
	senderPeerID peerID           // unique system id
	ibGID        net.IP           // gid of ib_device port
	ibMAC        net.HardwareAddr // mac of ib_device port
	ipAreaOffset uint16           // offset to IP address info area

	// Optional SMC-D info
	smcdGID  uint64 // ISM GID of requestor
	reserved [32]byte

	// IP/prefix info
	prefix          net.IP // subnet mask (rather prefix)
	prefixLen       uint8  // number of significant bits in mask
	reserved2       [2]byte
	ipv6PrefixesCnt uint8 // number of IPv6 prefixes in prefix array
	ipv6Prefixes    []ipv6Prefix

	trailer
}

// String converts the CLC Proposal message to a string
func (p *proposal) String() string {
	if p == nil {
		return "n/a"
	}

	// smc-d info
	smcdInfo := ""
	if p.ipAreaOffset == smcdIPAreaOffset {
		smcdInfo = fmt.Sprintf("SMC-D GID: %d, ", p.smcdGID)
	}

	// ipv6 prefixes
	ipv6Prefixes := ""
	for _, prefix := range p.ipv6Prefixes {
		ipv6Prefixes += fmt.Sprintf(", IPv6 Prefix: %s", prefix)
	}

	proposalFmt := "%s, Peer ID: %s, SMC-R GID: %s, RoCE MAC: %s, " +
		"IP Area Offset: %d, %sIPv4 Prefix: %s/%d, " +
		"IPv6 Prefix Count: %d%s, Trailer: %s"
	return fmt.Sprintf(proposalFmt, p.header.String(), p.senderPeerID,
		p.ibGID, p.ibMAC, p.ipAreaOffset, smcdInfo, p.prefix,
		p.prefixLen, p.ipv6PrefixesCnt, ipv6Prefixes, p.trailer)
}

// Reserved converts the CLC Proposal message to a string including reserved
// message fields
func (p *proposal) Reserved() string {
	if p == nil {
		return "n/a"
	}

	// smc-d info
	smcdInfo := ""
	if p.ipAreaOffset == smcdIPAreaOffset {
		smcdInfo = fmt.Sprintf("SMC-D GID: %d, Reserved: %#x, ",
			p.smcdGID, p.reserved)
	}
	// ipv6 prefixes
	ipv6Prefixes := ""
	for _, prefix := range p.ipv6Prefixes {
		ipv6Prefixes += fmt.Sprintf(", IPv6 Prefix: %s", prefix)
	}

	proposalFmt := "%s, Peer ID: %s, SMC-R GID: %s, RoCE MAC: %s, " +
		"IP Area Offset: %d, %sIPv4 Prefix: %s/%d, Reserved: %#x, " +
		"IPv6 Prefix Count: %d%s, Trailer: %s"
	return fmt.Sprintf(proposalFmt, p.header.Reserved(), p.senderPeerID,
		p.ibGID, p.ibMAC, p.ipAreaOffset, smcdInfo, p.prefix,
		p.prefixLen, p.reserved2, p.ipv6PrefixesCnt, ipv6Prefixes,
		p.trailer)
}

// Parse parses the CLC Proposal message in buf
func (p *proposal) Parse(buf []byte) {
	// save raw message bytes
	p.raw.Parse(buf)

	// parse CLC header
	p.header.Parse(buf)

	// check if message is long enough
	if p.Length < proposalLen {
		log.Println("Error parsing CLC Proposal: message too short")
		errDump(buf[:p.Length])
		return
	}

	// skip clc header
	skip := HeaderLen

	// sender peer ID
	copy(p.senderPeerID[:], buf[skip:skip+peerIDLen])
	skip += peerIDLen

	// ib GID is an IPv6 address
	p.ibGID = make(net.IP, net.IPv6len)
	copy(p.ibGID[:], buf[skip:skip+net.IPv6len])
	skip += net.IPv6len

	// ib MAC is a 6 byte MAC address
	p.ibMAC = make(net.HardwareAddr, 6)
	copy(p.ibMAC[:], buf[skip:skip+6])
	skip += 6

	// offset to ip area
	p.ipAreaOffset = binary.BigEndian.Uint16(buf[skip : skip+2])
	skip += 2

	// Optional SMC-D info
	if p.ipAreaOffset == smcdIPAreaOffset {
		// smcd GID
		p.smcdGID = binary.BigEndian.Uint64(buf[skip : skip+8])
		skip += 8

		// reserved
		copy(p.reserved[:], buf[skip:skip+32])
		skip += 32
	} else {
		skip += int(p.ipAreaOffset)
	}

	// make sure we do not read outside the message
	if int(p.Length)-skip < net.IPv4len+1+2+1+trailerLen {
		log.Println("Error parsing CLC Proposal: " +
			"IP Area Offset too big")
		errDump(buf[:p.Length])
		return
	}

	// IP/prefix is an IPv4 address
	p.prefix = make(net.IP, net.IPv4len)
	copy(p.prefix[:], buf[skip:skip+net.IPv4len])
	skip += net.IPv4len

	// prefix length
	p.prefixLen = uint8(buf[skip])
	skip++

	// reserved
	copy(p.reserved2[:], buf[skip:skip+2])
	skip += 2

	// ipv6 prefix count
	p.ipv6PrefixesCnt = uint8(buf[skip])

	// parse ipv6 prefixes
	for i := uint8(0); i < p.ipv6PrefixesCnt; i++ {
		// skip prefix count or last prefix length
		skip++

		// make sure we are still inside the clc message
		if int(p.Length)-skip < ipv6PrefixLen+trailerLen {
			log.Println("Error parsing CLC Proposal: " +
				"IPv6 prefix count too big")
			errDump(buf[:p.Length])
			break
		}
		// create new ipv6 prefix entry
		ip6prefix := ipv6Prefix{}

		// parse prefix and fill prefix entry
		ip6prefix.prefix = make(net.IP, net.IPv6len)
		copy(ip6prefix.prefix[:], buf[skip:skip+net.IPv6len])
		skip += net.IPv6len

		// parse prefix length and fill prefix entry
		ip6prefix.prefixLen = uint8(buf[skip])

		// add to ipv6 prefixes
		p.ipv6Prefixes = append(p.ipv6Prefixes, ip6prefix)
	}

	// save trailer
	p.trailer.Parse(buf)
}
