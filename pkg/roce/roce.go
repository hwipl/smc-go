package roce

import "github.com/hwipl/smc-go/pkg/llc"

const (
	// RoCEv1EtherType is the Ethernet Type used for RoCEv1
	RoCEv1EtherType = 0x8915

	// RoCEv2UDPPort is the UDP port number used for RoCEv2
	RoCEv2UDPPort = 4791
)

// RoCE stores a roce message
type RoCE struct {
	Type string
	GRH  *GRH
	BTH  *BTH
	LLC  llc.Message
	ICRC []byte
}

// ParseRoCEv1 parses the RoCEv1 packet in buffer
func ParseRoCEv1(buffer []byte) *RoCE {
	var roce RoCE

	// set type string
	roce.Type = "RoCEv1"

	// Global Routing Header (GRH) is 40 bytes (it's an IPv6 header)
	roce.GRH = ParseGRH(buffer[:GRHLen])
	buffer = buffer[GRHLen:]

	// Base Transport Header (BTH) is 12 bytes
	roce.BTH = ParseBTH(buffer[:BTHLen])
	buffer = buffer[BTHLen:]

	// invariant CRC (ICRC) is 4 bytes at the end of the message
	roce.ICRC = buffer[len(buffer)-4:]
	buffer = buffer[:len(buffer)-4]

	// parse payload
	roce.LLC = llc.ParseLLC(buffer)

	return &roce
}

// ParseRoCEv2 parses the RoCEv2 packet in buffer
func ParseRoCEv2(buffer []byte) *RoCE {
	var roce RoCE

	// set type string
	roce.Type = "RoCEv2"

	// Base Transport Header (BTH) is 12 bytes
	roce.BTH = ParseBTH(buffer[:BTHLen])
	buffer = buffer[BTHLen:]

	// invariant CRC (ICRC) is 4 bytes at the end of the message
	roce.ICRC = buffer[len(buffer)-4:]
	buffer = buffer[:len(buffer)-4]

	// parse payload
	roce.LLC = llc.ParseLLC(buffer)

	return &roce
}
