package roce

import (
	"encoding/hex"
	"fmt"

	"github.com/gopacket/gopacket"
	"github.com/gopacket/gopacket/layers"
)

const (
	// internal message type
	typeGRH = 0x102

	// GRHLen is the length of the GRH
	GRHLen = 40
)

// GRH stores an ib global routing header
type GRH layers.IPv6

// Parse fills the grh fields from the global routing header in buffer
func (g *GRH) Parse(buffer []byte) {
	// parse buffer as IPv6 packet
	ipv6Packet := gopacket.NewPacket(buffer, layers.LayerTypeIPv6,
		gopacket.Default)
	ipv6Layer := ipv6Packet.Layer(layers.LayerTypeIPv6)
	if ipv6Layer != nil {
		ipv6, _ := ipv6Layer.(*layers.IPv6)
		*g = GRH(*ipv6)
	}

}

// String converts the global routing header to a string
func (g *GRH) String() string {
	// rename nextHeader to "BTH" if it is correct (== 0x1B)
	nextHeader := "BTH"
	if g.NextHeader != BTHNextHeader {
		// otherwise just use what gopacket thinks it is
		nextHeader = g.NextHeader.String()
	}
	iFmt := "GRH: Version: %d, Traffic Class: %d, " +
		"Flow Label: %d, Length: %d, Next Header: %s, " +
		"Hop Limit: %d, Source: %s, Destination: %s\n"
	return fmt.Sprintf(iFmt, g.Version, g.TrafficClass, g.FlowLabel,
		g.Length, nextHeader, g.HopLimit, g.SrcIP, g.DstIP)
}

// Reserved converts the global routing header to a string including reserved
// fields
func (g *GRH) Reserved() string {
	return g.String()
}

// Hex returns a hex dump string of the message
func (g *GRH) Hex() string {
	return hex.Dump(g.Contents)
}

// GetType returns the type of the message
func (g *GRH) GetType() int {
	return typeGRH
}

// ParseGRH parses the global routing header in buffer
func ParseGRH(buffer []byte) *GRH {
	var g GRH
	g.Parse(buffer)
	return &g
}
