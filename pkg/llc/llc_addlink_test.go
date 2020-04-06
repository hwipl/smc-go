package llc

import (
	"encoding/hex"
	"log"
	"strings"
	"testing"
)

func TestAddLink(t *testing.T) {
	var want, got string

	// create bytes of a message
	msgHex := "00000000  02 2c 00 00 98 03 9b ab  " +
		"cd ef 00 00 fe 80 00 00  |.,..............|\n" +
		"00000010  00 00 00 00 9a 03 9b ff  " +
		"fe ab cd ef 00 00 00 00  |................|\n" +
		"00000020  00 00 00 00 00 00 00 00  " +
		"00 00 00 00              |............|\n"
	msg := "02 2c 00 00 98 03 9b ab  cd ef 00 00 fe 80 00 00" +
		"00 00 00 00 9a 03 9b ff  fe ab cd ef 00 00 00 00" +
		"00 00 00 00 00 00 00 00  00 00 00 00"
	bytes, err := hex.DecodeString(strings.Join(strings.Fields(msg), ""))
	if err != nil {
		log.Fatal(err)
	}

	// parse message
	a := parseAddLink(bytes)

	// test String()
	want = "LLC Add Link: Type: 2, Length: 44, " +
		"Reason Code: 0 (unknown), Reply: false, Rejection: false, " +
		"Sender MAC: 98:03:9b:ab:cd:ef, " +
		"Sender GID: fe80::9a03:9bff:feab:cdef, " +
		"Sender QP: 0, Link: 0, MTU: 0 (reserved), Initial PSN: 0\n"
	got = a.String()
	if got != want {
		t.Errorf("a.String() = %s; want %s", got, want)
	}

	// test Reserved()
	want = "LLC Add Link: Type: 2, Length: 44, Reserved: 0x0, " +
		"Reason Code: 0 (unknown), Reply: false, Rejection: false, " +
		"Reserved: 0x0, Sender MAC: 98:03:9b:ab:cd:ef, " +
		"Sender GID: fe80::9a03:9bff:feab:cdef, " +
		"Sender QP: 0, Link: 0, Reserved: 0x0, MTU: 0 (reserved), " +
		"Initial PSN: 0, Reserved: 0x00000000000000000000\n"
	got = a.Reserved()
	if got != want {
		t.Errorf("a.Reserved() = %s; want %s", got, want)
	}

	// test Hex()
	want = msgHex
	got = a.Hex()
	if got != want {
		t.Errorf("a.Hex() = %s; want %s", got, want)
	}

	// test GetType()
	wantType := typeAddLink
	gotType := a.GetType()
	if gotType != wantType {
		t.Errorf("a.GetType() = %d; want %d", gotType, wantType)
	}
}
