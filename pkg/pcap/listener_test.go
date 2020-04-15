package pcap

import (
	"io/ioutil"
	"log"
	"net"
	"os"
	"testing"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcapgo"
)

type testHandler struct {
	packetHandled bool
	timerHandled  bool
}

func (th *testHandler) HandlePacket(packet gopacket.Packet) {
	th.packetHandled = true
}

func (th *testHandler) HandleTimer() {
	th.timerHandled = true
}

func TestListenerPcap(t *testing.T) {
	// prepare creation of packet
	opts := gopacket.SerializeOptions{
		FixLengths:       true,
		ComputeChecksums: true,
	}

	// create ethernet header
	srcMAC, err := net.ParseMAC("00:00:5e:00:53:01")
	if err != nil {
		log.Fatal(err)
	}
	dstMAC := srcMAC
	eth := layers.Ethernet{
		SrcMAC: srcMAC,
		DstMAC: dstMAC,
	}

	// serialize packet to buffer
	buf := gopacket.NewSerializeBuffer()
	err = gopacket.SerializeLayers(buf, opts, &eth)
	if err != nil {
		log.Fatal(err)
	}
	packet := buf.Bytes()

	// create temporary pcap file
	tmpFile, err := ioutil.TempFile("", "listener.pcap")
	if err != nil {
		log.Fatal(err)
	}
	defer os.Remove(tmpFile.Name())

	// write packets of fake tcp connection to pcap file
	w := pcapgo.NewWriter(tmpFile)
	w.WriteFileHeader(65536, layers.LinkTypeEthernet)
	w.WritePacket(gopacket.CaptureInfo{
		CaptureLength: len(packet),
		Length:        len(packet),
	}, packet)
	tmpFile.Close()

	// prepare listener
	var testHandler testHandler
	var listener Listener
	listener.PacketHandler = &testHandler
	listener.File = tmpFile.Name()

	// let listener handle the packet
	listener.Prepare()
	listener.Loop()

	// check results
	want := true
	got := testHandler.packetHandled
	if got != want {
		t.Errorf("got = %t; want %t", got, want)
	}
}
