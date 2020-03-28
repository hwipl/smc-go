package clc

import (
	"encoding/binary"
	"fmt"
	"net"
)

const (
	// PeerIDLen is the length of the peer ID in bytes
	PeerIDLen = 8
)

// PeerID stores a SMC peer ID
type PeerID [PeerIDLen]byte

// String converts the peer ID to a string
func (p PeerID) String() string {
	instance := binary.BigEndian.Uint16(p[:2])
	roceMAC := net.HardwareAddr(p[2:8])
	return fmt.Sprintf("%d@%s", instance, roceMAC)
}
