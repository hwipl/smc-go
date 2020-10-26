package clc

import (
	"encoding/binary"
	"fmt"
	"log"
)

// clc operating system types
const (
	ZOS   = 1
	Linux = 2
	AIX   = 3
)

// OSType is the operating system type
type OSType uint8

// String converts OSType to a string
func (o OSType) String() string {
	var os string
	switch o {
	case ZOS:
		os = "z/OS"
	case Linux:
		os = "Linux"
	case AIX:
		os = "AIX"
	default:
		os = "unknown"
	}
	return fmt.Sprintf("%d (%s)", o, os)
}

// Decline stores a SMCv2 CLC Decline message
type DeclineV2 struct {
	Raw
	Header
	SenderPeerID  PeerID        // sender peer id
	PeerDiagnosis PeerDiagnosis // diagnosis information
	OSType        OSType        // OS type (4 bits)
	reserved      [4]byte       // first byte contains OS type
	Trailer
}

// String converts the SMCv2 CLC Decline message to a string
func (d *DeclineV2) String() string {
	if d == nil {
		return "n/a"
	}

	declineFmt := "%s, Peer ID: %s, Peer Diagnosis: %s, " +
		"OS Type: %s, Trailer: %s"
	return fmt.Sprintf(declineFmt, d.Header.String(),
		d.SenderPeerID, d.PeerDiagnosis, d.OSType, d.Trailer)
}

// Reserved converts the SMCv2 CLC Decline message to a string including
// reserved message fields
func (d *DeclineV2) Reserved() string {
	if d == nil {
		return "n/a"
	}

	declineFmt := "%s, Peer ID: %s, Peer Diagnosis: %s, " +
		"OS Type: %s, Reserved: %#x, Trailer: %s"
	return fmt.Sprintf(declineFmt, d.Header.Reserved(),
		d.SenderPeerID, d.PeerDiagnosis, d.OSType,
		d.reserved, d.Trailer)
}

// Parse parses the SMCv2 CLC Decline message in buf
func (d *DeclineV2) Parse(buf []byte) {
	// save raw message bytes
	d.Raw.Parse(buf)

	// parse CLC header
	d.Header.Parse(buf)

	// check if message is long enough
	if d.Length < DeclineLen {
		log.Println("Error parsing CLC Decline: message too short")
		errDump(buf[:d.Length])
		return
	}

	// skip clc header
	buf = buf[HeaderLen:]

	// sender peer ID
	copy(d.SenderPeerID[:], buf[:PeerIDLen])
	buf = buf[PeerIDLen:]

	// peer diagnosis
	d.PeerDiagnosis = PeerDiagnosis(binary.BigEndian.Uint32(buf[:4]))
	buf = buf[4:]

	// reserved
	copy(d.reserved[:], buf[:4])

	// os type (4 highest bits of first byte of reserved)
	d.OSType = OSType(buf[0] >> 4)
	d.reserved[0] &= 0b00001111
	buf = buf[4:]

	// save trailer
	d.Trailer.Parse(buf)
}
