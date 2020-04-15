package pcap

import (
	"fmt"
	"log"
	"time"

	"github.com/google/gopacket"
	"github.com/google/gopacket/pcap"
)

// Listener is a pcap listener that reads packets from a file or device and
// calls Handlers for packets and timer events
type Listener struct {
	pcapHandle *pcap.Handle

	PacketHandler PacketHandler

	Timer        time.Duration
	TimerHandler TimerHandler

	File    string
	Device  string
	Promisc bool
	Snaplen int
	Timeout time.Duration
	Filter  string
	MaxPkts int
	MaxTime time.Duration
}

// getFirstPcapInterface sets the first network interface found by pcap
func (p *Listener) getFirstPcapInterface() {
	ifs, err := pcap.FindAllDevs()
	if err != nil {
		log.Fatal(err)
	}
	if len(ifs) > 0 {
		p.Device = ifs[0].Name
		return
	}
	log.Fatal("No network interface found")
}

// Prepare prepares the pcap listener for the listen function
func (p *Listener) Prepare() {
	// open pcap handle
	var pcapErr error
	var startText string
	if p.File == "" {
		// set pcap timeout
		timeout := pcap.BlockForever
		if p.Timeout > 0 {
			timeout = p.Timeout
		}

		// set interface
		if p.Device == "" {
			p.getFirstPcapInterface()
		}

		// open device
		p.pcapHandle, pcapErr = pcap.OpenLive(p.Device,
			int32(p.Snaplen), p.Promisc, timeout)
		startText = fmt.Sprintf("Listening on interface %s:\n",
			p.Device)
	} else {
		// open pcap file
		p.pcapHandle, pcapErr = pcap.OpenOffline(p.File)
		startText = fmt.Sprintf("Reading packets from file %s:\n",
			p.File)
	}
	if pcapErr != nil {
		log.Fatal(pcapErr)
	}
	if p.Filter != "" {
		if err := p.pcapHandle.SetBPFFilter(p.Filter); err != nil {
			log.Fatal(pcapErr)
		}
	}
	log.Printf(startText)
}

// Loop implements the listen loop for the listen function
func (p *Listener) Loop() {
	defer p.pcapHandle.Close()

	// Use the handle as a packet source to process all packets
	packetSource := gopacket.NewPacketSource(p.pcapHandle,
		p.pcapHandle.LinkType())
	packets := packetSource.Packets()

	// setup timer
	ticker := time.Tick(p.Timer)

	// set stop time if configured
	stop := make(<-chan time.Time)
	if p.MaxTime > 0 {
		stop = time.After(p.MaxTime)
	}

	// handle packets and timer events
	count := 0
	for {
		select {
		case packet := <-packets:
			if packet == nil {
				return
			}
			p.PacketHandler.HandlePacket(packet)
			count++
			if p.MaxPkts > 0 && count == p.MaxPkts {
				return
			}
		case <-ticker:
			p.TimerHandler.HandleTimer()
		case <-stop:
			return
		}
	}

}
