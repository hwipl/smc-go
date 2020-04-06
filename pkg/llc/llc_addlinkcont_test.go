package llc

import (
	"encoding/hex"
	"log"
	"strings"
	"testing"
)

func TestAddLinkCont(t *testing.T) {
	var want, got string

	// create bytes of a fake message
	msgHex := "00000000  03 2c 00 80 01 01 00 00  " +
		"00 00 00 01 00 00 00 02  |.,..............|\n" +
		"00000010  00 00 00 00 00 00 00 03  " +
		"00 00 00 04 00 00 00 05  |................|\n" +
		"00000020  00 00 00 00 00 00 00 06  " +
		"00 00 00 00              |............|\n"
	msg := "03 2c 00 80 01 01 00 00  00 00 00 01 00 00 00 02" +
		"00 00 00 00 00 00 00 03 00 00 00 04 00 00 00 05" +
		"00 00 00 00 00 00 00 06 00 00 00 00"
	bytes, err := hex.DecodeString(strings.Join(strings.Fields(msg), ""))
	if err != nil {
		log.Fatal(err)
	}

	// parse message
	a := parseAddLinkCont(bytes)

	// test String()
	want = "LLC Add Link Continuation: Type: 3, Length: 44, " +
		"Reply: true, Link: 1, Number of RTokens: 1, " +
		"RKey Pair 2: [Reference RKey: 4, New RKey: 5, " +
		"New Virtual Address: 0x6]\n"
	got = a.String()
	if got != want {
		t.Errorf("a.String() = %s; want %s", got, want)
	}

	// test Reserved()
	want = "LLC Add Link Continuation: Type: 3, Length: 44, " +
		"Reserved: 0x0, Reply: true, Reserved: 0x0, Link: 1, " +
		"Number of RTokens: 1, Reserved: 0x0000, " +
		"RKey Pair 2: [Reference RKey: 4, New RKey: 5, " +
		"New Virtual Address: 0x6], Reserved: 0x00000000\n"
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
	wantType := typeAddLinkCont
	gotType := a.GetType()
	if gotType != wantType {
		t.Errorf("a.GetType() = %d; want %d", gotType, wantType)
	}

}
