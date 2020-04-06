package llc

import (
	"encoding/hex"
	"log"
	"strings"
	"testing"
)

func TestConfirmRKey(t *testing.T) {
	var want, got string

	// create bytes of a message
	msgHex := "00000000  06 2c 00 00 00 00 00 15  " +
		"5d 00 00 00 0e ef 4d 00  |.,......].....M.|\n" +
		"00000010  00 00 00 00 00 00 00 00  " +
		"00 00 00 00 00 00 00 00  |................|\n" +
		"00000020  00 00 00 00 00 00 00 00  " +
		"00 00 00 00              |............|\n"
	msg := "06 2c 00 00 00 00 00 15  5d 00 00 00 0e ef 4d 00" +
		"00 00 00 00 00 00 00 00  00 00 00 00 00 00 00 00" +
		"00 00 00 00 00 00 00 00  00 00 00 00"
	bytes, err := hex.DecodeString(strings.Join(strings.Fields(msg), ""))
	if err != nil {
		log.Fatal(err)
	}

	// parse message
	c := ParseConfirmRKey(bytes)

	// test String()
	want = "LLC Confirm RKey: Type: 6, Length: 44, Reply: false, " +
		"Negative Response: false, Configuration Retry: false, " +
		"Number of Tokens: 0, This RKey: 5469, " +
		"This VAddr: 0xeef4d0000, Other Link RMB 1: " +
		"[Link: 0, RKey: 0, Virtual Address: 0x0], " +
		"Other Link RMB 2: [Link: 0, RKey: 0, " +
		"Virtual Address: 0x0]\n"
	got = c.String()
	if got != want {
		t.Errorf("c.String() = %s; want %s", got, want)
	}

	// test Reserved()
	want = "LLC Confirm RKey: Type: 6, Length: 44, Reserved: 0x0, " +
		"Reply: false, Reserved: 0x0, Negative Response: false, " +
		"Configuration Retry: false, Reserved: 0x0, " +
		"Number of Tokens: 0, This RKey: 5469, " +
		"This VAddr: 0xeef4d0000, Other Link RMB 1: " +
		"[Link: 0, RKey: 0, Virtual Address: 0x0], " +
		"Other Link RMB 2: [Link: 0, RKey: 0, " +
		"Virtual Address: 0x0], Reserved: 0x0\n"
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
	wantType := TypeConfirmRKey
	gotType := c.GetType()
	if gotType != wantType {
		t.Errorf("c.GetType() = %d; want %d", gotType, wantType)
	}

}
