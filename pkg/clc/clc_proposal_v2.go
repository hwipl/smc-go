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
