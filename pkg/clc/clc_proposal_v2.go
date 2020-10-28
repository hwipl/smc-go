package clc

import (
	"encoding/binary"
	"fmt"
	"log"
	"net"
)

// SMCv2 constants
const (
	// ProposalV2Len is the minimum length of a Proposal v2 message
	ProposalV2Len = 84
	// EIDLen is the length of an EID
	EIDLen = 32
	// ProposalV2ExtLen is the minimum Proposal v2 Extension length
	ProposalV2ExtLen = 40
	// SMCDv2ExtLen is the minimum SMC-D v2 Extension length
	SMCDv2ExtLen = 48
)

// GIDEntry stores a SMC-D GID entry consisting of GID and VCHID
type GIDEntry struct {
	GID   uint64
	VCHID uint16
}

// ProposalV2 stores a SMCv2 CLC Proposal message
type ProposalV2 struct {
	Raw
	Header
	SenderPeerID PeerID // unique system id

	// SMC-R GID info
	IBGID        net.IP           // gid of ib_device port
	IBMAC        net.HardwareAddr // mac of ib_device port
	IPAreaOffset uint16           // offset to IP address info area

	// SMC-D GID info
	SMCDGID     uint64 // ISM GID of requestor
	ISMv2VCHID  uint16 // ISMv2 VCHID
	SMCv2Offset uint16 // SMCv2 Extension Offset
	reserved    [28]byte

	// Optional IP/Prefix info
	Prefix          net.IP // subnet mask (rather prefix)
	PrefixLen       uint8  // number of significant bits in mask
	reserved2       [2]byte
	IPv6PrefixesCnt uint8 // number of IPv6 prefixes in prefix array
	IPv6Prefixes    []IPv6Prefix

	// CLC Proposal Message V2 Extension
	EIDNumber uint8 // Number of EIDs in the EID Array Area
	GIDNumber uint8 // Number of GIDs in the ISMv2 GID Array Area
	reserved3 byte
	Release   uint8 // SMCv2 Release number (4 bits)
	reserved4 byte  // 3 bits
	SEIDInd   uint8 // SEID indicator (1 bit): 0 not present, 1 present
	reserved5 [2]byte
	SMCDv2Off uint16 // SMC-Dv2 Extension Offset (if present)
	reserved6 [32]byte
	EIDArea   [8][EIDLen]byte // stores 0-8 EIDs, see EIDNumber

	// Optional SMC-Dv2 Extension
	SEID      [EIDLen]byte
	reserved7 [16]byte
	GIDArea   [8]GIDEntry // stores 0-8 GIDs/VCHIDs, see GIDNumber

	Trailer
}

// ipV4String converts the ipv4 info to a string
func (p *ProposalV2) ipV4String() string {
	return fmt.Sprintf("IPv4 Prefix: %s/%d, ", p.Prefix, p.PrefixLen)
}

// ipV6String converts the ipv6 info to a string
func (p *ProposalV2) ipV6String() string {
	ipInfo := fmt.Sprintf("IPv6 Prefix Count: %d", p.IPv6PrefixesCnt)
	for _, prefix := range p.IPv6Prefixes {
		ipInfo += fmt.Sprintf(", IPv6 Prefix: %s", prefix)
	}
	return ipInfo
}

// ipInfoString converts the optional IP/Prefix information to a string
func (p *ProposalV2) ipInfoString() string {
	if p.Path == SMCTypeN {
		return ""
	}

	return p.ipV4String() + p.ipV6String()
}

// propV2ExtEIDString converts the Proposal v2 Extension EID Area to a string
func (p *ProposalV2) propV2ExtEIDString() string {
	eidArea := ""
	for i, eid := range p.EIDArea {
		if i >= int(p.EIDNumber) {
			break
		}
		if i > 0 {
			eidArea += ", "
		}
		eidArea += fmt.Sprintf("EID %d: %s", i, eid)
	}
	return eidArea
}

// propV2ExtString converts the Proposal v2 Extension to a string
func (p *ProposalV2) propV2ExtString() string {
	if p.Pathv2 == SMCTypeN {
		return ""
	}

	// EID area
	eidArea := p.propV2ExtEIDString()

	extFmt := "EID Number: %d, GID Number: %d, Release: %d, " +
		"SEID Indicator: %d, SMC-Dv2 Extension Offset: %d, " +
		"EID Area: [%s]"
	return fmt.Sprintf(extFmt, p.EIDNumber, p.GIDNumber, p.Release,
		p.SEIDInd, p.SMCDv2Off, eidArea)
}

// smcdV2ExtString converts the SMC-D v2 Extension GID Area to a string
func (p *ProposalV2) smcdV2ExtGIDString() string {
	gidArea := ""
	for i, gid := range p.GIDArea {
		if i >= int(p.GIDNumber) {
			break
		}
		if i > 0 {
			gidArea += ", "
		}
		gidArea += fmt.Sprintf("GID %d: %d, VCHID %d: %d", i, gid.GID,
			i, gid.VCHID)
	}
	return gidArea
}

// smcdV2ExtString converts the SMC-D v2 Extension to a string
func (p *ProposalV2) smcdV2ExtString() string {
	if p.Pathv2 != SMCTypeD && p.Pathv2 != SMCTypeB {
		return ""
	}

	// GID area
	gidArea := p.smcdV2ExtGIDString()

	extFmt := "SEID: %s, GID Area: [%s]"
	return fmt.Sprintf(extFmt, p.SEID, gidArea)
}

// String converts the CLC Proposal message to a string
func (p *ProposalV2) String() string {
	if p == nil {
		return "n/a"
	}

	// optional ip/prefix info
	ipInfo := p.ipInfoString()
	if ipInfo != "" {
		ipInfo = ", " + ipInfo
	}

	// clc proposal message v2 extension
	propV2Ext := p.propV2ExtString()
	if propV2Ext != "" {
		propV2Ext = ", " + propV2Ext
	}

	// smc-d v2 extension
	smcdV2Ext := p.smcdV2ExtString()
	if smcdV2Ext != "" {
		smcdV2Ext = ", " + smcdV2Ext
	}

	proposalFmt := "%s, Peer ID: %s, SMC-R GID: %s, RoCE MAC: %s, " +
		"IP Area Offset: %d, SMC-D GID: %d, ISMv2 VCHID: %d, " +
		"SMCv2 Extension Offset: %d%s%s%s, Trailer: %s"
	return fmt.Sprintf(proposalFmt, p.Header.String(), p.SenderPeerID,
		p.IBGID, p.IBMAC, p.IPAreaOffset, p.SMCDGID, p.ISMv2VCHID,
		p.SMCv2Offset, ipInfo, propV2Ext, smcdV2Ext, p.Trailer)
}

// ipInfoReserved converts the optional IP/Prefix information to a string
// including reserved message fields
func (p *ProposalV2) ipInfoReserved() string {
	if p.Path == SMCTypeN {
		return ""
	}

	return p.ipV4String() + fmt.Sprintf("Reserved: %#x, ", p.reserved2) +
		p.ipV6String()
}

// propV2ExtReserved converts the Proposal v2 Extension to a string including
// reserved message fields
func (p *ProposalV2) propV2ExtReserved() string {
	if p.Pathv2 == SMCTypeN {
		return ""
	}

	// EID area
	eidArea := p.propV2ExtEIDString()

	extFmt := "EID Number: %d, GID Number: %d, Reserved: %#x, " +
		"Release: %d, Reserved: %#x, SEID Indicator: %d, " +
		"Reserved: %#x, SMC-Dv2 Extension Offset: %d, Reserved: %#x, " +
		"EID Area: [%s]"
	return fmt.Sprintf(extFmt, p.EIDNumber, p.GIDNumber, p.reserved3,
		p.Release, p.reserved4, p.SEIDInd, p.reserved5, p.SMCDv2Off,
		p.reserved6, eidArea)
}

// smcdV2ExtReserved converts the SMC-D v2 Extension to a string including
// reserved message fields
func (p *ProposalV2) smcdV2ExtReserved() string {
	if p.Pathv2 != SMCTypeD && p.Pathv2 != SMCTypeB {
		return ""
	}

	// GID area
	gidArea := p.smcdV2ExtGIDString()

	extFmt := "SEID: %s, Reserved: %#v, GID Area: [%s]"
	return fmt.Sprintf(extFmt, p.SEID, p.reserved7, gidArea)
}

// Reserved converts the CLC Proposal message to a string including reserved
// message fields
func (p *ProposalV2) Reserved() string {
	if p == nil {
		return "n/a"
	}

	// optional ip/prefix info
	ipInfo := p.ipInfoReserved()
	if ipInfo != "" {
		ipInfo = ", " + ipInfo
	}

	// clc proposal message v2 extension
	propV2Ext := p.propV2ExtReserved()
	if propV2Ext != "" {
		propV2Ext = ", " + propV2Ext
	}

	// smc-d v2 extension
	smcdV2Ext := p.smcdV2ExtReserved()
	if smcdV2Ext != "" {
		smcdV2Ext = ", " + smcdV2Ext
	}
	proposalFmt := "%s, Peer ID: %s, SMC-R GID: %s, RoCE MAC: %s, " +
		"IP Area Offset: %d, SMC-D GID: %d, ISMv2 VCHID: %d, " +
		"SMCv2 Extension Offset: %d, Reserved: %#x%s%s%s, Trailer: %s"
	return fmt.Sprintf(proposalFmt, p.Header.Reserved(), p.SenderPeerID,
		p.IBGID, p.IBMAC, p.IPAreaOffset, p.SMCDGID, p.ISMv2VCHID,
		p.SMCv2Offset, p.reserved, ipInfo, propV2Ext, smcdV2Ext,
		p.Trailer)
}

// Parse parses the SMCv2 CLC Proposal message in buf
func (p *ProposalV2) Parse(buf []byte) {
	// save raw message bytes
	p.Raw.Parse(buf)

	// parse CLC header
	p.Header.Parse(buf)

	// check if message is long enough
	if p.Length < ProposalV2Len {
		log.Println("Error parsing CLC Proposal v2: message too short")
		errDump(buf[:p.Length])
		return
	}

	// skip clc header
	skip := HeaderLen

	// sender peer ID
	copy(p.SenderPeerID[:], buf[skip:skip+PeerIDLen])
	skip += PeerIDLen

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

	// smcd GID
	p.SMCDGID = binary.BigEndian.Uint64(buf[skip : skip+8])
	skip += 8

	// ism v2 vchid
	p.ISMv2VCHID = binary.BigEndian.Uint16(buf[skip : skip+2])
	skip += 2

	// smc v2 extension offset
	p.SMCv2Offset = binary.BigEndian.Uint16(buf[skip : skip+2])
	skip += 2

	// reserved
	copy(p.reserved[:], buf[skip:skip+28])
	skip += 28

	// parse optional ip/prefix info
	if p.Path != SMCTypeN {
		// make sure we do not read outside the message
		if int(p.Length)-skip < net.IPv4len+1+2+1+TrailerLen {
			log.Println("Error parsing CLC Proposal v2: " +
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
		skip++

		// parse ipv6 prefixes
		for i := uint8(0); i < p.IPv6PrefixesCnt; i++ {
			// make sure we are still inside the clc message
			if int(p.Length)-skip < IPv6PrefixLen+TrailerLen {
				log.Println("Error parsing CLC Proposal v2: " +
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
			skip++

			// add to ipv6 prefixes
			p.IPv6Prefixes = append(p.IPv6Prefixes, ip6prefix)
		}
	}

	// parse proposal message v2 extension
	if p.Pathv2 != SMCTypeN {
		// make sure we do not read outside the message
		if int(p.Length)-skip < ProposalV2ExtLen+TrailerLen {
			log.Println("Error parsing CLC Proposal v2: " +
				"Not enough space for Proposal v2 Extension")
			errDump(buf[:p.Length])
			return
		}

		// number of EIDs in EID Area
		p.EIDNumber = buf[skip]
		skip++

		//number of GIDs in ISMv2 GID Array Area
		p.GIDNumber = buf[skip]
		skip++

		// reserved
		p.reserved3 = buf[skip]
		skip++

		// Release number (4 bits)
		p.Release = buf[skip] >> 4
		// reserved (3 bits)
		p.reserved4 = (buf[skip] & 0b00001110) >> 1
		// SEID indicator (1 bit)
		p.SEIDInd = buf[skip] & 0b00000001
		skip++

		// reserved
		copy(p.reserved5[:], buf[skip:skip+2])
		skip += 2

		// smcd v2 extension offset
		p.SMCDv2Off = binary.BigEndian.Uint16(buf[skip : skip+2])
		skip += 2

		// reserved
		copy(p.reserved6[:], buf[skip:skip+32])
		skip += 32

		// parse EIDs in EID Area
		for i := uint8(0); i < p.EIDNumber; i++ {
			// make sure we are still inside the clc message
			if int(p.Length)-skip < EIDLen+TrailerLen {
				log.Println("Error parsing CLC Proposal v2: " +
					"EID number too big")
				errDump(buf[:p.Length])
				break
			}

			// parse EID
			copy(p.EIDArea[i][:], buf[skip:skip+EIDLen])
			skip += EIDLen
		}
	}

	// parse optional smcd v2 extension
	if p.Pathv2 == SMCTypeD || p.Pathv2 == SMCTypeB {
		// make sure we do not read outside the message
		if int(p.Length)-skip < SMCDv2ExtLen+TrailerLen {
			log.Println("Error parsing CLC Proposal v2: " +
				"Not enough space for SMC-D v2 Extension")
			errDump(buf[:p.Length])
			return
		}

		// SEID
		copy(p.SEID[:], buf[skip:skip+32])
		skip += 32

		// reserved
		copy(p.reserved7[:], buf[skip:skip+16])
		skip += 16

		// parse GIDs in GID Area
		for i := uint8(0); i < p.GIDNumber; i++ {
			// make sure we are still inside the clc message
			if int(p.Length)-skip < 8+2+TrailerLen {
				log.Println("Error parsing CLC Proposal v2: " +
					"GID number too big")
				errDump(buf[:p.Length])
				break
			}

			// parse GID
			p.GIDArea[i].GID = binary.BigEndian.Uint64(
				buf[skip : skip+8])
			skip += 8

			// parse VCHID
			p.GIDArea[i].VCHID = binary.BigEndian.Uint16(
				buf[skip : skip+2])
			skip += 2
		}
	}

	// save trailer
	p.Trailer.Parse(buf)
}
