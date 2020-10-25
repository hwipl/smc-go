package clc

import (
	"encoding/binary"
	"fmt"
	"log"
)

const (
	// DeclineLen is the length of a clc decline message in bytes
	DeclineLen = 28
)

// decline diagnosis codes (linux)
const (
	DeclineMem        = 0x01010000 // insufficient memory resources
	DeclineTimeoutCL  = 0x02010000 // timeout w4 QP confirm link
	DeclineTimeoutAL  = 0x02020000 // timeout w4 QP add link
	DeclineCnfErr     = 0x03000000 // configuration error
	DeclinePeerNoSMC  = 0x03010000 // peer did not indicate SMC
	DeclineIPSEC      = 0x03020000 // IPsec usage
	DeclineNoSMCDev   = 0x03030000 // no SMC device found (R or D)
	DeclineNoSMCDDev  = 0x03030001 // no SMC-D device found
	DeclineNoSMCRDev  = 0x03030002 // no SMC-R device found
	DeclineSMCDNoTalk = 0x03030003 // SMC-D dev can't talk to peer
	DeclineModeUnsupp = 0x03040000 // smc modes do not match (R or D)
	DeclineRMBEEyeC   = 0x03050000 // peer has eyecatcher in RMBE
	DeclineOptUnsupp  = 0x03060000 // fastopen sockopt not supported
	DeclineDiffPrefix = 0x03070000 // IP prefix / subnet mismatch
	DeclineGetVLANErr = 0x03080000 // err to get vlan id of ip device
	DeclineISMVLANErr = 0x03090000 // err to reg vlan id on ism dev
	DeclineSyncErr    = 0x04000000 // synchronization error
	DeclinePeerDecl   = 0x05000000 // peer declined during handshake
	DeclineInterr     = 0x09990000 // internal error
	DeclineErrRTok    = 0x09990001 // rtoken handling failed
	DeclineErrRdyLnk  = 0x09990002 // ib ready link failed
	DeclineErrRegRMB  = 0x09990003 // reg rmb failed
)

// clc operating system types
const (
	ZOS   = 1
	Linux = 2
	AIX   = 3
)

// PeerDiagnosis stores the decline diagnosis code in a decline message
type PeerDiagnosis uint32

// String converts the peerDiagnosis to a string
func (p PeerDiagnosis) String() string {
	// parse peer diagnosis code
	var diag string
	switch p {
	case DeclineMem:
		diag = "insufficient memory resources"
	case DeclineTimeoutCL:
		diag = "timeout w4 QP confirm link"
	case DeclineTimeoutAL:
		diag = "timeout w4 QP add link"
	case DeclineCnfErr:
		diag = "configuration error"
	case DeclinePeerNoSMC:
		diag = "peer did not indicate SMC"
	case DeclineIPSEC:
		diag = "IPsec usage"
	case DeclineNoSMCDev:
		diag = "no SMC device found (R or D)"
	case DeclineNoSMCDDev:
		diag = "no SMC-D device found"
	case DeclineNoSMCRDev:
		diag = "no SMC-R device found"
	case DeclineSMCDNoTalk:
		diag = "SMC-D dev can't talk to peer"
	case DeclineModeUnsupp:
		diag = "smc modes do not match (R or D)"
	case DeclineRMBEEyeC:
		diag = "peer has eyecatcher in RMBE"
	case DeclineOptUnsupp:
		diag = "fastopen sockopt not supported"
	case DeclineDiffPrefix:
		diag = "IP prefix / subnet mismatch"
	case DeclineGetVLANErr:
		diag = "err to get vlan id of ip device"
	case DeclineISMVLANErr:
		diag = "err to reg vlan id on ism dev"
	case DeclineSyncErr:
		diag = "synchronization error"
	case DeclinePeerDecl:
		diag = "peer declined during handshake"
	case DeclineInterr:
		diag = "internal error"
	case DeclineErrRTok:
		diag = "rtoken handling failed"
	case DeclineErrRdyLnk:
		diag = "ib ready link failed"
	case DeclineErrRegRMB:
		diag = "reg rmb failed"
	default:
		diag = "Unknown"
	}
	return fmt.Sprintf("%#x (%s)", uint32(p), diag)
}

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

// Decline stores a CLC Decline message
type Decline struct {
	Raw
	Header
	SenderPeerID  PeerID        // sender peer id
	PeerDiagnosis PeerDiagnosis // diagnosis information
	reserved      [4]byte
	Trailer
}

// String converts the CLC Decline message to a string
func (d *Decline) String() string {
	if d == nil {
		return "n/a"
	}

	declineFmt := "%s, Peer ID: %s, Peer Diagnosis: %s, Trailer: %s"
	return fmt.Sprintf(declineFmt, d.Header.String(), d.SenderPeerID,
		d.PeerDiagnosis, d.Trailer)
}

// Reserved converts the CLC Decline message to a string including reserved
// message fields
func (d *Decline) Reserved() string {
	if d == nil {
		return "n/a"
	}

	declineFmt := "%s, Peer ID: %s, Peer Diagnosis: %s, Reserved: %#x, " +
		"Trailer: %s"
	return fmt.Sprintf(declineFmt, d.Header.Reserved(), d.SenderPeerID,
		d.PeerDiagnosis, d.reserved, d.Trailer)
}

// Parse parses the CLC Decline in buf
func (d *Decline) Parse(buf []byte) {
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
	buf = buf[4:]

	// save trailer
	d.Trailer.Parse(buf)
}
