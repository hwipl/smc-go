package pcap

// Handler is the interface for handlers called by Listener
type Handler interface {
	HandlePacket(gopacket.Packet)
	HandleTimer()
}
