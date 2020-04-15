package pcap

import "github.com/google/gopacket"

// PacketHandler is the interface called by Listener for handling packets
type PacketHandler interface {
	Handle(gopacket.Packet)
}

// TimerHandler is the interface called by Listener for timer events
type TimerHandler interface {
	Handle()
}
