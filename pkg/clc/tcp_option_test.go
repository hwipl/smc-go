package clc

import (
	"testing"

	"github.com/gopacket/gopacket/layers"
)

func TestCheckSMCOption(t *testing.T) {
	var want, got bool
	var tcp layers.TCP

	// test packet without option (empty tcp layer)
	want = false
	got = CheckSMCOption(&tcp)
	if got != want {
		t.Errorf("got = %t; want %t", got, want)
	}

	// append the smc option and test again
	opt := layers.TCPOption{
		OptionType:   254,
		OptionLength: 6,
		OptionData:   SMCREyecatcher,
	}
	tcp.Options = append(tcp.Options, opt)
	want = true
	got = CheckSMCOption(&tcp)
	if got != want {
		t.Errorf("got = %t; want %t", got, want)
	}
}
