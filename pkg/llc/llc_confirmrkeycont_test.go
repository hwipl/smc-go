package llc

import (
	"encoding/hex"
	"log"
	"strings"
	"testing"
)

func TestConfirmRKeyCont(t *testing.T) {
	var want, got string

	// create bytes of a fake message
	msgHex := "00000000  08 2c 00 a0 01 01 00 00  " +
		"00 02 00 00 00 00 00 00  |.,..............|\n" +
		"00000010  00 03 04 00 00 00 05 00  " +
		"00 00 00 00 00 00 06 07  |................|\n" +
		"00000020  00 00 00 08 00 00 00 00  " +
		"00 00 00 09 00           |.............|\n"
	msg := "08 2c 00 a0 01 01 00 00  00 02 00 00 00 00 00 00" +
		"00 03 04 00 00 00 05 00  00 00 00 00 00 00 06 07" +
		"00 00 00 08 00 00 00 00  00 00 00 09 00"
	bytes, err := hex.DecodeString(strings.Join(strings.Fields(msg), ""))
	if err != nil {
		log.Fatal(err)
	}

	// parse message
	c := parseConfirmRKeyCont(bytes)

	// test String()
	want = "LLC Confirm RKey Continuation: Type: 8, Length: 44, " +
		"Reply: true, Negative Response: true, " +
		"Number of Tokens: 1, Other Link RMB 1: [Link: 1, " +
		"RKey: 2, Virtual Address: 0x3], Other Link RMB 2: " +
		"[Link: 4, RKey: 5, Virtual Address: 0x6], " +
		"Other Link RMB 3: [Link: 7, RKey: 8, " +
		"Virtual Address: 0x9]\n"
	got = c.String()
	if got != want {
		t.Errorf("c.String() = %s; want %s", got, want)
	}

	// test Reserved()
	want = "LLC Confirm RKey Continuation: Type: 8, Length: 44, " +
		"Reserved: 0x0, Reply: true, Reserved: 0x0, " +
		"Negative Response: true, Reserved: 0x0, " +
		"Number of Tokens: 1, Other Link RMB 1: [Link: 1, " +
		"RKey: 2, Virtual Address: 0x3], Other Link RMB 2: " +
		"[Link: 4, RKey: 5, Virtual Address: 0x6], " +
		"Other Link RMB 3: [Link: 7, RKey: 8, " +
		"Virtual Address: 0x9], Reserved: 0x0\n"
	got = c.Reserved()
	if got != want {
		t.Errorf("c.Reserved() = %s; want %s", got, want)
	}

	// test Hex()
	want = msgHex
	got = c.Hex()
	if got != want {
		t.Errorf("c.Hex() = %s; want %s", got, want)
	}

	// test GetType()
	wantType := typeConfirmRKeyCont
	gotType := c.GetType()
	if gotType != wantType {
		t.Errorf("c.GetType() = %d; want %d", gotType, wantType)
	}

}
