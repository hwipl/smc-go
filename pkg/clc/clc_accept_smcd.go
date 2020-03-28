package clc

import (
	"encoding/binary"
	"fmt"
	"log"
)

const (
	AcceptSMCDLen = 48
)

// AcceptSMCD stores a CLC SMC-D Accept message
type AcceptSMCD struct {
	Raw
	Header
	GID       uint64   // Sender GID
	Token     uint64   // DMB token
	DMBEIdx   uint8    // DMBE index
	DMBESize  RMBESize // 4 bits buf size (compressed)
	reserved  byte     // 4 bits reserved
	reserved2 [2]byte
	LinkID    uint32 // Link identifier
	reserved3 [12]byte
	trailer
}

// String converts the CLC SMC-D Accept message to a string
func (ac *AcceptSMCD) String() string {
	if ac == nil {
		return "n/a"
	}

	acFmt := "%s, SMC-D GID: %d, SMC-D Token: %d, DMBE Index: %d, " +
		"DMBE Size: %s, Link ID: %d, Trailer: %s"
	return fmt.Sprintf(acFmt, ac.Header.String(), ac.GID, ac.Token,
		ac.DMBEIdx, ac.DMBESize, ac.LinkID, ac.trailer)
}

// Reserved converts the CLC SMC-D Accept message to a string including
// reserved message fields
func (ac *AcceptSMCD) Reserved() string {
	if ac == nil {
		return "n/a"
	}

	acFmt := "%s, SMC-D GID: %d, SMC-D Token: %d, DMBE Index: %d, " +
		"DMBE Size: %s, Reserved: %#x, Reserved: %#x, " +
		"Link ID: %d, Reserved: %#x, Trailer: %s"
	return fmt.Sprintf(acFmt, ac.Header.Reserved(), ac.GID,
		ac.Token, ac.DMBEIdx, ac.DMBESize, ac.reserved,
		ac.reserved2, ac.LinkID, ac.reserved3, ac.trailer)
}

// Parse parses the SMC-D Accept message in buf
func (ac *AcceptSMCD) Parse(buf []byte) {
	// save raw message bytes
	ac.Raw.Parse(buf)

	// parse clc header
	ac.Header.Parse(buf)

	// check if message is long enough
	if ac.Length < AcceptSMCDLen {
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

	// reserved
	copy(ac.reserved3[:], buf[:12])
	buf = buf[12:]

	// save trailer
	ac.trailer.Parse(buf)
}
