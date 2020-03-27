package clc

import (
	"encoding/hex"
	"log"
	"testing"
)

func TestParseSMCRConfirm(t *testing.T) {
	// prepare message
	msgBytes := "e2d4c3d903004410b1a098039babcdef" +
		"fe800000000000009a039bfffeabcdef" +
		"98039babcdef0000e50000187f010000" +
		"0006230000000000f0a40000000d89a4" +
		"e2d4c3d9"
	msg, err := hex.DecodeString(msgBytes)
	if err != nil {
		log.Fatal(err)
	}

	// parse message
	ac, acLen := NewMessage(msg)
	ac.Parse(msg)

	// check message length
	if acLen != 68 {
		t.Errorf("acLen = %d; want %d", acLen, 68)
	}

	// check output message without reserved fields
	hdr := "Confirm: Eyecatcher: SMC-R, Type: 3 (Confirm), " +
		"Length: 68, Version: 1, Flag: 0, Path: SMC-R, "
	mid := "Peer ID: 45472@98:03:9b:ab:cd:ef, " +
		"SMC-R GID: fe80::9a03:9bff:feab:cdef, " +
		"RoCE MAC: 98:03:9b:ab:cd:ef, QP Number: 229, " +
		"RMB RKey: 6271, RMBE Index: 1, RMBE Alert Token: 6, " +
		"RMBE Size: 2 (65536), QP MTU: 3 (1024), " +
		"RMB Virtual Address: 0xf0a40000, " +
		"Packet Sequence Number: 887204"
	trl := ", Trailer: SMC-R"
	want := hdr + mid + trl
	got := ac.String()
	if got != want {
		t.Errorf("ac.String() = %s; want %s", got, want)
	}

	// check output message with reserved fields
	hdr = "Confirm: Eyecatcher: SMC-R, Type: 3 (Confirm), " +
		"Length: 68, Version: 1, Flag: 0, Reserved: 0x0, " +
		"Path: SMC-R, "
	mid = "Peer ID: 45472@98:03:9b:ab:cd:ef, " +
		"SMC-R GID: fe80::9a03:9bff:feab:cdef, " +
		"RoCE MAC: 98:03:9b:ab:cd:ef, QP Number: 229, " +
		"RMB RKey: 6271, RMBE Index: 1, RMBE Alert Token: 6, " +
		"RMBE Size: 2 (65536), QP MTU: 3 (1024), Reserved: 0x0, " +
		"RMB Virtual Address: 0xf0a40000, Reserved: 0x0, " +
		"Packet Sequence Number: 887204"
	trl = ", Trailer: SMC-R"
	want = hdr + mid + trl
	got = ac.Reserved()
	if got != want {
		t.Errorf("ac.Reserved() = %s; want %s", got, want)
	}
}
