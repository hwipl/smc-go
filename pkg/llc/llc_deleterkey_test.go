package llc

import (
	"encoding/hex"
	"log"
	"strings"
	"testing"
)

func TestDeleteRKey(t *testing.T) {
	var want, got string

	// create bytes of a message
	msgHex := "00000000  09 2c 00 00 01 00 00 00  " +
		"00 00 15 5d 00 00 00 00  |.,.........]....|\n" +
		"00000010  00 00 00 00 00 00 00 00  " +
		"00 00 00 00 00 00 00 00  |................|\n" +
		"00000020  00 00 00 00 00 00 00 00  " +
		"00 00 00 00              |............|\n"
	msg := "09 2c 00 00 01 00 00 00  00 00 15 5d 00 00 00 00" +
		"00 00 00 00 00 00 00 00  00 00 00 00 00 00 00 00" +
		"00 00 00 00 00 00 00 00  00 00 00 00"
	bytes, err := hex.DecodeString(strings.Join(strings.Fields(msg), ""))
	if err != nil {
		log.Fatal(err)
	}

	// parse message
	d := parseDeleteRKey(bytes)

	// test String()
	want = "LLC Delete RKey: Type: 9, Length: 44, Reply: false, " +
		"Negative Response: false, Count: 1, Error Mask: 0b0, " +
		"RKey 0: 5469, RKey 1: 0, RKey 2: 0, RKey 3: 0, " +
		"RKey 4: 0, RKey 5: 0, RKey 6: 0, RKey 7: 0\n"
	got = d.String()
	if got != want {
		t.Errorf("d.String() = %s; want %s", got, want)
	}

	// test Reserved()
	want = "LLC Delete RKey: Type: 9, Length: 44, Reserved: 0x0, " +
		"Reply: false, Reserved: 0x0, Negative Response: false, " +
		"Reserved: 0x0, Count: 1, Error Mask: 0b0, " +
		"Reserved: 0x0000, RKey 0: 5469, RKey 1: 0, RKey 2: 0, " +
		"RKey 3: 0, RKey 4: 0, RKey 5: 0, RKey 6: 0, RKey 7: 0, " +
		"Reserved: 0x00000000\n"
	got = d.Reserved()
	if got != want {
		t.Errorf("d.Reserved() = %s; want %s", got, want)
	}

	// test Hex()
	want = msgHex
	got = d.Hex()
	if got != want {
		t.Errorf("d.Hex() = %s; want %s", got, want)
	}

	// test GetType()
	wantType := typeDeleteRKey
	gotType := d.GetType()
	if gotType != wantType {
		t.Errorf("d.GetType() = %d; want %d", gotType, wantType)
	}

}
