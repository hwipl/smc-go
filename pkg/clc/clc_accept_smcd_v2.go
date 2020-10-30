package clc

import (
	"fmt"
)

const (
	// AcceptSMCDv2Len is the miminum length of a SMCv2 SMC-D Accept msg
	AcceptSMCDv2Len = 78
	// AcceptSMCDv2FCELen is the length of a SMCv2 SMC-D Accept with FCE
	AcceptSMCDv2FCELen = 114
)

// AcceptSMCDv2 stores a SMCv2 CLC SMC-D Accept message
type AcceptSMCDv2 struct {
	Raw
	Header
	GID        uint64   // Sender GID
	Token      uint64   // DMB token
	DMBEIdx    uint8    // DMBE index
	DMBESize   RMBESize // 4 bits buf size (compressed)
	reserved   byte     // 4 bits reserved
	reserved2  [2]byte
	LinkID     uint32 // Link identifier
	ISMv2VCHID uint16 // ISMv2 VCHID
	EID        EID    // EID
	reserved3  [8]byte

	// First Contact Extension; only present if first contact flag set
	reserved4 byte
	OSType    OSType // 4 bits
	Release   uint8  // 4 bits
	reserved5 [2]byte
	Hostname  EID // hostname has same format as EID

	Trailer
}

// fceString converts the FCE in the SMCv2 CLC SMC-D Accept message to a string
func (ac *AcceptSMCDv2) fceString() string {
	if ac.Length < AcceptSMCDv2FCELen {
		return ""
	}
	fceFmt := ", OS Type: %s, Release: %d, Hostname: %s"
	return fmt.Sprintf(fceFmt, ac.OSType, ac.Release, &ac.Hostname)
}

// String converts the SMCv2 CLC SMC-D Accept message to a string
func (ac *AcceptSMCDv2) String() string {
	if ac == nil {
		return "n/a"
	}

	acFmt := "%s, SMC-D GID: %d, SMC-D Token: %d, DMBE Index: %d, " +
		"DMBE Size: %s, Link ID: %d, ISMv2 VCHID: %d, EID: %s%s, " +
		"Trailer: %s"
	return fmt.Sprintf(acFmt, ac.Header.String(), ac.GID, ac.Token,
		ac.DMBEIdx, ac.DMBESize, ac.LinkID, ac.ISMv2VCHID, &ac.EID,
		ac.fceString(), ac.Trailer)
}

// fceReserved converts the FCE in the SMCv2 CLC SMC-D Accept messate to a
// string including reserved message fields
func (ac *AcceptSMCDv2) fceReserved() string {
	if ac.Length < AcceptSMCDv2FCELen {
		return ""
	}
	fceFmt := ", Reserved: %#x, OS Type: %s, Release: %d, " +
		"Reserved: %#x, Hostname: %s"
	return fmt.Sprintf(fceFmt, ac.reserved4, ac.OSType, ac.Release,
		ac.reserved5, &ac.Hostname)
}

// Reserved converts the SMCv2 CLC SMC-D Accept message to a string including
// reserved message fields
func (ac *AcceptSMCDv2) Reserved() string {
	if ac == nil {
		return "n/a"
	}

	acFmt := "%s, SMC-D GID: %d, SMC-D Token: %d, DMBE Index: %d, " +
		"DMBE Size: %s, Reserved: %#x, Reserved: %#x, " +
		"Link ID: %d, ISMv2 VCHID: %d, EID: %s, Reserved: %#x%s, " +
		"Trailer: %s"
	return fmt.Sprintf(acFmt, ac.Header.Reserved(), ac.GID,
		ac.Token, ac.DMBEIdx, ac.DMBESize, ac.reserved,
		ac.reserved2, ac.LinkID, ac.ISMv2VCHID, &ac.EID,
		ac.reserved3, ac.fceReserved(), ac.Trailer)
}
