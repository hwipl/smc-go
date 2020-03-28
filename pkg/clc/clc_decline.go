package clc

import (
	"encoding/binary"
	"fmt"
	"log"
)

const (
	// declineLen is the length of a clc decline message in bytes
	declineLen = 28

	// decline diagnosis codes (linux)
	declineMem        = 0x01010000 // insufficient memory resources
	declineTimeoutCL  = 0x02010000 // timeout w4 QP confirm link
	declineTimeoutAL  = 0x02020000 // timeout w4 QP add link
	declineCnfErr     = 0x03000000 // configuration error
	declinePeerNoSMC  = 0x03010000 // peer did not indicate SMC
	declineIPSEC      = 0x03020000 // IPsec usage
	declineNoSMCDev   = 0x03030000 // no SMC device found (R or D)
	declineNoSMCDDev  = 0x03030001 // no SMC-D device found
	declineNoSMCRDev  = 0x03030002 // no SMC-R device found
	declineSMCDNoTalk = 0x03030003 // SMC-D dev can't talk to peer
	declineModeUnsupp = 0x03040000 // smc modes do not match (R or D)
	declineRMBEEyeC   = 0x03050000 // peer has eyecatcher in RMBE
	declineOptUnsupp  = 0x03060000 // fastopen sockopt not supported
	declineDiffPrefix = 0x03070000 // IP prefix / subnet mismatch
	declineGetVLANErr = 0x03080000 // err to get vlan id of ip device
	declineISMVLANErr = 0x03090000 // err to reg vlan id on ism dev
	declineSyncErr    = 0x04000000 // synchronization error
	declinePeerDecl   = 0x05000000 // peer declined during handshake
	declineInterr     = 0x09990000 // internal error
	declineErrRTok    = 0x09990001 // rtoken handling failed
	declineErrRdyLnk  = 0x09990002 // ib ready link failed
	declineErrRegRMB  = 0x09990003 // reg rmb failed
)

// peerDiagnosis stores the decline diagnosis code in a decline message
type peerDiagnosis uint32

// String converts the peerDiagnosis to a string
func (p peerDiagnosis) String() string {
	// parse peer diagnosis code
	var diag string
	switch p {
	case declineMem:
		diag = "insufficient memory resources"
	case declineTimeoutCL:
		diag = "timeout w4 QP confirm link"
	case declineTimeoutAL:
		diag = "timeout w4 QP add link"
	case declineCnfErr:
		diag = "configuration error"
	case declinePeerNoSMC:
		diag = "peer did not indicate SMC"
	case declineIPSEC:
		diag = "IPsec usage"
	case declineNoSMCDev:
		diag = "no SMC device found (R or D)"
	case declineNoSMCDDev:
		diag = "no SMC-D device found"
	case declineNoSMCRDev:
		diag = "no SMC-R device found"
	case declineSMCDNoTalk:
		diag = "SMC-D dev can't talk to peer"
	case declineModeUnsupp:
		diag = "smc modes do not match (R or D)"
	case declineRMBEEyeC:
		diag = "peer has eyecatcher in RMBE"
	case declineOptUnsupp:
		diag = "fastopen sockopt not supported"
	case declineDiffPrefix:
		diag = "IP prefix / subnet mismatch"
	case declineGetVLANErr:
		diag = "err to get vlan id of ip device"
	case declineISMVLANErr:
		diag = "err to reg vlan id on ism dev"
	case declineSyncErr:
		diag = "synchronization error"
	case declinePeerDecl:
		diag = "peer declined during handshake"
	case declineInterr:
		diag = "internal error"
	case declineErrRTok:
		diag = "rtoken handling failed"
	case declineErrRdyLnk:
		diag = "ib ready link failed"
	case declineErrRegRMB:
		diag = "reg rmb failed"
	default:
		diag = "Unknown"
	}
	return fmt.Sprintf("%#x (%s)", uint32(p), diag)
}

// Decline stores a CLC Decline message
type Decline struct {
	raw
	header
	SenderPeerID  peerID        // sender peer id
	PeerDiagnosis peerDiagnosis // diagnosis information
	reserved      [4]byte
	trailer
}

// String converts the CLC Decline message to a string
func (d *Decline) String() string {
	if d == nil {
		return "n/a"
	}

	declineFmt := "%s, Peer ID: %s, Peer Diagnosis: %s, Trailer: %s"
	return fmt.Sprintf(declineFmt, d.header.String(), d.SenderPeerID,
		d.PeerDiagnosis, d.trailer)
}

// Reserved converts the CLC Decline message to a string including reserved
// message fields
func (d *Decline) Reserved() string {
	if d == nil {
		return "n/a"
	}

	declineFmt := "%s, Peer ID: %s, Peer Diagnosis: %s, Reserved: %#x, " +
		"Trailer: %s"
	return fmt.Sprintf(declineFmt, d.header.Reserved(), d.SenderPeerID,
		d.PeerDiagnosis, d.reserved, d.trailer)
}

// Parse parses the CLC Decline in buf
func (d *Decline) Parse(buf []byte) {
	// save raw message bytes
	d.raw.Parse(buf)

	// parse CLC header
	d.header.Parse(buf)

	// check if message is long enough
	if d.Length < declineLen {
		log.Println("Error parsing CLC Decline: message too short")
		errDump(buf[:d.Length])
		return
	}

	// skip clc header
	buf = buf[HeaderLen:]

	// sender peer ID
	copy(d.SenderPeerID[:], buf[:peerIDLen])
	buf = buf[peerIDLen:]

	// peer diagnosis
	d.PeerDiagnosis = peerDiagnosis(binary.BigEndian.Uint32(buf[:4]))
	buf = buf[4:]

	// reserved
	copy(d.reserved[:], buf[:4])
	buf = buf[4:]

	// save trailer
	d.trailer.Parse(buf)
}
