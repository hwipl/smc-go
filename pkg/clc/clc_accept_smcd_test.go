package clc

import (
	"encoding/hex"
	"log"
	"testing"
)

func TestParseSMCDAccept(t *testing.T) {
	// prepare message
	msgBytes := "e2d4c3c4020030110123456789abcdef" +
		"0123456789abcdefff100000ffffffff" +
		"000000000000000000000000e2d4c3c4"
	msg, err := hex.DecodeString(msgBytes)
	if err != nil {
		log.Fatal(err)
	}

	// parse message
	clc, clcLen := NewMessage(msg)
	clc.Parse(msg)

	// check message length
	if clcLen != 48 {
		t.Errorf("clcLen = %d; want %d", clcLen, 48)
	}

	// check output message without reserved fields
	hdr := "Accept: Eyecatcher: SMC-D, Type: 2 (Accept), " +
		"Length: 48, Version: 1, First Contact: 0, Path: SMC-D, "
	mid := "SMC-D GID: 81985529216486895, " +
		"SMC-D Token: 81985529216486895, " +
		"DMBE Index: 255, DMBE Size: 1 (32768), Link ID: 4294967295"
	trl := ", Trailer: SMC-D"
	want := hdr + mid + trl
	got := clc.String()
	if got != want {
		t.Errorf("clc.String() = %s; want %s", got, want)
	}

	// check output message with reserved fields
	hdr = "Accept: Eyecatcher: SMC-D, Type: 2 (Accept), " +
		"Length: 48, Version: 1, First Contact: 0, Reserved: 0x0, " +
		"Path: SMC-D, "
	mid = "SMC-D GID: 81985529216486895, " +
		"SMC-D Token: 81985529216486895, " +
		"DMBE Index: 255, DMBE Size: 1 (32768), Reserved: 0x0, " +
		"Reserved: 0x0000, Link ID: 4294967295, " +
		"Reserved: 0x000000000000000000000000"
	trl = ", Trailer: SMC-D"
	want = hdr + mid + trl
	got = clc.Reserved()
	if got != want {
		t.Errorf("clc.Reserved() = %s; want %s", got, want)
	}
}
