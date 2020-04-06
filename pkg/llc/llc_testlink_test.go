package llc

import (
	"encoding/hex"
	"log"
	"strings"
	"testing"
)

func TestTestLink(t *testing.T) {
	var want, got string

	// create bytes of a fake message
	msgHex := "00000000  07 2c 00 80 00 00 00 00  " +
		"00 00 00 00 00 00 00 00  |.,..............|\n" +
		"00000010  00 00 00 00 00             " +
		"                       |.....|\n"
	msg := "07 2c 00 80 00 00 00 00  00 00 00 00 00 00 00 00" +
		"00 00 00 00 00"
	bytes, err := hex.DecodeString(strings.Join(strings.Fields(msg), ""))
	if err != nil {
		log.Fatal(err)
	}

	// parse message
	tl := parseTestLink(bytes)

	// test String()
	want = "LLC Test Link: Type 7, Length: 44, Reply: true, " +
		"User Data: 0x00000000000000000000000000000000\n"
	got = tl.String()
	if got != want {
		t.Errorf("tl.String() = %s; want %s", got, want)
	}

	// test Reserved()
	want = "LLC Test Link: Type 7, Length: 44, Reserved: 0x0, " +
		"Reply: true, Reserved: 0x0, " +
		"User Data: 0x00000000000000000000000000000000, " +
		"Reserved: 0x00000000000000000000000000000000000" +
		"0000000000000\n"
	got = tl.Reserved()
	if got != want {
		t.Errorf("tl.Reserved() = %s; want %s", got, want)
	}

	// test Hex()
	want = msgHex
	got = tl.Hex()
	if got != want {
		t.Errorf("tl.Hex() = %s; want %s", got, want)
	}

	// test GetType()
	wantType := typeTestLink
	gotType := tl.GetType()
	if gotType != wantType {
		t.Errorf("tl.GetType() = %d; want %d", gotType, wantType)
	}

}
