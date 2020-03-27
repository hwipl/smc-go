package clc

import (
	"encoding/binary"
	"fmt"
	"net"
)

const (
	// peerIDLen is the length of the peer ID in bytes
	peerIDLen = 8
)

// peerID stores a SMC peer ID
type peerID [peerIDLen]byte

// String converts the peer ID to a string
func (p peerID) String() string {
	instance := binary.BigEndian.Uint16(p[:2])
	roceMAC := net.HardwareAddr(p[2:8])
	return fmt.Sprintf("%d@%s", instance, roceMAC)
}
