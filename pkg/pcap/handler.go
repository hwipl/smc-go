package pcap

import "github.com/google/gopacket"

// PacketHandler is the interface called by Listener for handling packets
type PacketHandler interface {
	HandlePacket(gopacket.Packet)
}

// TimerHandler is the interface called by Listener for timer events
type TimerHandler interface {
	HandleTimer()
}
