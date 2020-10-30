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
