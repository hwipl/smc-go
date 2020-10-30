package clc

import (
	"encoding/binary"
	"fmt"
	"log"
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

// Parse parses the SMCv2 CLC SMC-D Accept message in buf
func (ac *AcceptSMCDv2) Parse(buf []byte) {
	// save raw message bytes
	ac.Raw.Parse(buf)

	// parse clc header
	ac.Header.Parse(buf)

	// check if message is long enough
	if ac.Length < AcceptSMCDv2Len || len(buf) < AcceptSMCDv2Len {
		err := "Error parsing SMC-Dv2 CLC Accept: message too short"
		if ac.Type == TypeConfirm {
			err = "Error parsing SMC-Dv2 CLC Confirm: " +
				"message too short"
		}
		log.Println(err)
		errDump(buf[:ac.Length])
		return
	}

	// skip clc header
	buf = buf[HeaderLen:]

	// smcd GID
	ac.GID = binary.BigEndian.Uint64(buf[:8])
	buf = buf[8:]

	// smcd Token
	ac.Token = binary.BigEndian.Uint64(buf[:8])
	buf = buf[8:]

	// dmbe index
	ac.DMBEIdx = uint8(buf[0])
	buf = buf[1:]

	// 1 byte bitfield: dmbe size (4 bits), reserved (4 bits)
	ac.DMBESize = RMBESize((uint8(buf[0]) & 0b11110000) >> 4)
	ac.reserved = buf[0] & 0b00001111
	buf = buf[1:]

	// reserved
	copy(ac.reserved2[:], buf[:2])
	buf = buf[2:]

	// link id
	ac.LinkID = binary.BigEndian.Uint32(buf[:4])
	buf = buf[4:]

	// ISMv2 VCHID
	ac.ISMv2VCHID = binary.BigEndian.Uint16(buf[:2])
	buf = buf[2:]

	// EID
	copy(ac.EID[:], buf[:EIDLen])
	buf = buf[EIDLen:]

	// reserved
	copy(ac.reserved3[:], buf[:8])
	buf = buf[8:]

	// parse First Contact Extension (FCE) if present
	if ac.Length == AcceptSMCDv2FCELen {
		// reserved
		ac.reserved4 = buf[0]
		buf = buf[1:]

		// OS type (4 bits)
		ac.OSType = OSType(buf[0] >> 4)
		// Release (4 bits)
		ac.Release = buf[0] & 0b00001111
		buf = buf[1:]

		// reserved
		copy(ac.reserved5[:], buf[:2])
		buf = buf[2:]

		// hostname
		copy(ac.Hostname[:], buf[:EIDLen])
		buf = buf[EIDLen:]
	}

	// save trailer
	ac.Trailer.Parse(buf)
}
