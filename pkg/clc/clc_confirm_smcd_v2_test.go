package clc

import (
	"encoding/hex"
	"log"
	"testing"
)

func TestParseSMCDv2FCEConfirm(t *testing.T) {
	// prepare message
	msgBytes := "e2d4c3c4" + "03" + "0072" + "2" + "9" + "0123456789abcdef" +
		"0123456789abcdef" + "ff" + "1" + "0" + "0000" + "ffffffff" +
		"0123" +
		"546869734973534d4376324549443031" +
		"00000000000000000000000000000000" +
		"0000000000000000" +
		// fce
		"00" + "2" + "0" + "0000" +
		"546869734973486f73746e616d653031" +
		"00000000000000000000000000000000" +
		"e2d4c3c4"
	msg, err := hex.DecodeString(msgBytes)
	if err != nil {
		log.Fatal(err)
	}

	// parse message
	clc, clcLen := NewMessage(msg)
	clc.Parse(msg)

	// check message length
	if clcLen != AcceptSMCDv2FCELen {
		t.Errorf("clcLen = %d; want %d", clcLen, AcceptSMCDv2FCELen)
	}

	// check output message without reserved fields
	want := "Confirm: Eyecatcher: SMC-D, Type: 3 (Confirm), Length: 114, " +
		"Version: 2, First Contact: 1, Path: SMC-D, " +
		"SMC-D GID: 81985529216486895, " +
		"SMC-D Token: 81985529216486895, DMBE Index: 255, " +
		"DMBE Size: 1 (32768), Link ID: 4294967295, " +
		"ISMv2 VCHID: 291, EID: ThisIsSMCv2EID01, " +
		"OS Type: 2 (Linux), Release: 0, " +
		"Hostname: ThisIsHostname01, Trailer: SMC-D"
	got := clc.String()
	if got != want {
		t.Errorf("clc.String() = %s; want %s", got, want)
	}

	// check output message with reserved fields
	want = "Confirm: Eyecatcher: SMC-D, Type: 3 (Confirm), " +
		"Length: 114, Version: 2, First Contact: 1, Reserved: 0x0, " +
		"Path: SMC-D, SMC-D GID: 81985529216486895, " +
		"SMC-D Token: 81985529216486895, " +
		"DMBE Index: 255, DMBE Size: 1 (32768), Reserved: 0x0, " +
		"Reserved: 0x0000, Link ID: 4294967295, " +
		"ISMv2 VCHID: 291, EID: ThisIsSMCv2EID01, " +
		"Reserved: 0x0000000000000000, Reserved: 0x0, " +
		"OS Type: 2 (Linux), Release: 0, Reserved: 0x0000, " +
		"Hostname: ThisIsHostname01, Trailer: SMC-D"
	got = clc.Reserved()
	if got != want {
		t.Errorf("clc.Reserved() = %s; want %s", got, want)
	}
}
