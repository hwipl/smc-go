package clc

import (
	"encoding/hex"
	"log"
	"testing"
)

func TestParseSMCRAccept(t *testing.T) {
	// prepare message
	msgBytes := "e2d4c3d902004418b1a098039babcdef" +
		"fe800000000000009a039bfffeabcdef" +
		"98039babcdef0000e40000157d010000" +
		"0005230000000000f0a600000072f5fe" +
		"e2d4c3d9"
	msg, err := hex.DecodeString(msgBytes)
	if err != nil {
		log.Fatal(err)
	}

	// parse message
	clc, clcLen := NewMessage(msg)
	clc.Parse(msg)

	// check message length
	if clcLen != 68 {
		t.Errorf("clcLen = %d; want %d", clcLen, 68)
	}

	// check output message without reserved fields
	hdr := "Accept: Eyecatcher: SMC-R, Type: 2 (Accept), " +
		"Length: 68, Version: 1, First Contact: 1, Path: SMC-R, "
	mid := "Peer ID: 45472@98:03:9b:ab:cd:ef, " +
		"SMC-R GID: fe80::9a03:9bff:feab:cdef, " +
		"RoCE MAC: 98:03:9b:ab:cd:ef, QP Number: 228, " +
		"RMB RKey: 5501, RMBE Index: 1, RMBE Alert Token: 5, " +
		"RMBE Size: 2 (65536), QP MTU: 3 (1024), " +
		"RMB Virtual Address: 0xf0a60000, " +
		"Packet Sequence Number: 7534078"
	trl := ", Trailer: SMC-R"
	want := hdr + mid + trl
	got := clc.String()
	if got != want {
		t.Errorf("clc.String() = %s; want %s", got, want)
	}

	// check output message with reserved fields
	hdr = "Accept: Eyecatcher: SMC-R, Type: 2 (Accept), " +
		"Length: 68, Version: 1, First Contact: 1, Reserved: 0x0, " +
		"Path: SMC-R, "
	mid = "Peer ID: 45472@98:03:9b:ab:cd:ef, " +
		"SMC-R GID: fe80::9a03:9bff:feab:cdef, " +
		"RoCE MAC: 98:03:9b:ab:cd:ef, QP Number: 228, " +
		"RMB RKey: 5501, RMBE Index: 1, RMBE Alert Token: 5, " +
		"RMBE Size: 2 (65536), QP MTU: 3 (1024), Reserved: 0x0, " +
		"RMB Virtual Address: 0xf0a60000, Reserved: 0x0, " +
		"Packet Sequence Number: 7534078"
	trl = ", Trailer: SMC-R"
	want = hdr + mid + trl
	got = clc.Reserved()
	if got != want {
		t.Errorf("clc.Reserved() = %s; want %s", got, want)
	}
}
