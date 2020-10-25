package clc

import (
	"encoding/hex"
	"log"
	"testing"
)

func TestParseCLCDecline(t *testing.T) {
	// prepare decline message
	declineMsg := "e2d4c3d904001c102525252525252500" +
		"0303000000000000e2d4c3d9"
	msg, err := hex.DecodeString(declineMsg)
	if err != nil {
		log.Fatal(err)
	}

	// parse message
	decline, declineLen := NewMessage(msg)
	decline.Parse(msg)

	// check message length
	if declineLen != 28 {
		t.Errorf("declineLen = %d; want %d", declineLen, 28)
	}

	// check output message without reserved fields
	hdr := "Decline: Eyecatcher: SMC-R, Type: 4 (Decline), Length: 28, " +
		"Version: 1, Out of Sync: 0, Path: SMC-R, "
	mid := "Peer ID: 9509@25:25:25:25:25:00, " +
		"Peer Diagnosis: 0x3030000 (no SMC device found (R or D))"
	trl := ", Trailer: SMC-R"
	want := hdr + mid + trl
	got := decline.String()
	if got != want {
		t.Errorf("decline.String() = %s; want %s", got, want)
	}

	// check output message with reserved fields
	hdr = "Decline: Eyecatcher: SMC-R, Type: 4 (Decline), Length: 28, " +
		"Version: 1, Out of Sync: 0, Reserved: 0x0, Path: SMC-R, "
	mid = "Peer ID: 9509@25:25:25:25:25:00, " +
		"Peer Diagnosis: 0x3030000 (no SMC device found (R or D)), " +
		"Reserved: 0x00000000"
	trl = ", Trailer: SMC-R"
	want = hdr + mid + trl
	got = decline.Reserved()
	if got != want {
		t.Errorf("decline.Reserved() = %s; want %s", got, want)
	}
}

func TestParseCLCDeclineV2(t *testing.T) {
	// prepare decline message
	declineMsg := "e2d4c3d904001c202525252525252500" +
		"0303000020000000e2d4c3d9"
	msg, err := hex.DecodeString(declineMsg)
	if err != nil {
		log.Fatal(err)
	}

	// parse message
	decline, declineLen := NewMessage(msg)
	decline.Parse(msg)

	// check message length
	if declineLen != 28 {
		t.Errorf("declineLen = %d; want %d", declineLen, 28)
	}

	// check output message without reserved fields
	hdr := "Decline: Eyecatcher: SMC-R, Type: 4 (Decline), Length: 28, " +
		"Version: 2, Out of Sync: 0, Path: SMC-R, "
	mid := "Peer ID: 9509@25:25:25:25:25:00, " +
		"Peer Diagnosis: 0x3030000 (no SMC device found (R or D)), " +
		"OS Type: 2 (Linux)"
	trl := ", Trailer: SMC-R"
	want := hdr + mid + trl
	got := decline.String()
	if got != want {
		t.Errorf("decline.String() = %s; want %s", got, want)
	}

	// check output message with reserved fields
	hdr = "Decline: Eyecatcher: SMC-R, Type: 4 (Decline), Length: 28, " +
		"Version: 2, Out of Sync: 0, Reserved: 0x0, Path: SMC-R, "
	mid = "Peer ID: 9509@25:25:25:25:25:00, " +
		"Peer Diagnosis: 0x3030000 (no SMC device found (R or D)), " +
		"OS Type: 2 (Linux), Reserved: 0x00000000"
	trl = ", Trailer: SMC-R"
	want = hdr + mid + trl
	got = decline.Reserved()
	if got != want {
		t.Errorf("decline.Reserved() = %s; want %s", got, want)
	}
}
