package clc

import (
	"encoding/hex"
	"log"
	"testing"
)

func TestParseSMCDConfirm(t *testing.T) {
	// prepare message
	msgBytes := "e2d4c3c4030030110123456789abcdef" +
		"0123456789abcdefff100000ffffffff" +
		"000000000000000000000000e2d4c3c4"
	msg, err := hex.DecodeString(msgBytes)
	if err != nil {
		log.Fatal(err)
	}

	// parse message
	ac, acLen := NewMessage(msg)
	ac.Parse(msg)

	// check message length
	if acLen != 48 {
		t.Errorf("acLen = %d; want %d", acLen, 48)
	}

	// check output message without reserved fields
	hdr := "Confirm: Eyecatcher: SMC-D, Type: 3 (Confirm), " +
		"Length: 48, Version: 1, Flag: 0, Path: SMC-D, "
	mid := "SMC-D GID: 81985529216486895, " +
		"SMC-D Token: 81985529216486895, " +
		"DMBE Index: 255, DMBE Size: 1 (32768), Link ID: 4294967295"
	trl := ", Trailer: SMC-D"
	want := hdr + mid + trl
	got := ac.String()
	if got != want {
		t.Errorf("ac.String() = %s; want %s", got, want)
	}

	// check output message with reserved fields
	hdr = "Confirm: Eyecatcher: SMC-D, Type: 3 (Confirm), " +
		"Length: 48, Version: 1, Flag: 0, Reserved: 0x0, " +
		"Path: SMC-D, "
	mid = "SMC-D GID: 81985529216486895, " +
		"SMC-D Token: 81985529216486895, " +
		"DMBE Index: 255, DMBE Size: 1 (32768), Reserved: 0x0, " +
		"Reserved: 0x0000, Link ID: 4294967295, " +
		"Reserved: 0x000000000000000000000000"
	trl = ", Trailer: SMC-D"
	want = hdr + mid + trl
	got = ac.Reserved()
	if got != want {
		t.Errorf("ac.Reserved() = %s; want %s", got, want)
	}
}
