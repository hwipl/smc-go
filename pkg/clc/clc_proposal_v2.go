package clc

import (
	"fmt"
	"net"
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
	EIDArea   [8][32]byte // stores 0-8 EIDs, see EIDNumber

	// Optional SMC-Dv2 Extension
	SEID      [32]byte
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
		"SEID Indicator: %d, SMCDv2 Extension Offset: %d, " +
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
		gidArea += fmt.Sprintf("GID %d: %d, VCHID %d, %d", i, gid.GID,
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
		"IP Area Offset: %d, SMC-D GID: %d, ISMv2 VCHID: %d%s%s%s, " +
		"Trailer: %s"
	return fmt.Sprintf(proposalFmt, p.Header.String(), p.SenderPeerID,
		p.IBGID, p.IBMAC, p.IPAreaOffset, p.SMCDGID, p.ISMv2VCHID,
		ipInfo, propV2Ext, smcdV2Ext, p.Trailer)
}
