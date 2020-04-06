package llc

import (
	"encoding/hex"
	"log"
	"strings"
	"testing"
)

func TestConfirmLink(t *testing.T) {
	var want, got string

	// create bytes of a message
	msgHex := "00000000  01 2c 00 03 98 03 9b ab  " +
		"cd ef fe 80 00 00 00 00  |.,..............|\n" +
		"00000010  00 00 9a 03 9b ff fe ab  " +
		"cd ef 00 01 07 00 00 01  |................|\n" +
		"00000020  00 00 02 00 00 00 00 00  " +
		"00 00 00 00              |............|\n"
	msg := "01 2c 00 03 98 03 9b ab  cd ef fe 80 00 00 00 00" +
		"00 00 9a 03 9b ff fe ab  cd ef 00 01 07 00 00 01" +
		"00 00 02 00 00 00 00 00  00 00 00 00"
	bytes, err := hex.DecodeString(strings.Join(strings.Fields(msg), ""))
	if err != nil {
		log.Fatal(err)
	}

	// parse message
	c := ParseConfirm(bytes)

	// test String()
	want = "LLC Confirm Link: Type: 1, Length: 44, Reply: false, " +
		"Sender MAC: 98:03:9b:ab:cd:ef, " +
		"Sender GID: fe80::9a03:9bff:feab:cdef, " +
		"Sender QP: 263, Link: 0, Sender Link UserID: 65536, " +
		"Max Links: 2\n"
	got = c.String()
	if got != want {
		t.Errorf("c.String() = %s; want %s", got, want)
	}

	// test Reserved()
	want = "LLC Confirm Link: Type: 1, Length: 44, Reserved: 0x0, " +
		"Reply: false, Reserved: 0x3, " +
		"Sender MAC: 98:03:9b:ab:cd:ef, " +
		"Sender GID: fe80::9a03:9bff:feab:cdef, " +
		"Sender QP: 263, Link: 0, Sender Link UserID: 65536, " +
		"Max Links: 2, Reserved: 0x000000000000000000\n"
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
	wantType := typeConfirmLink
	gotType := c.GetType()
	if gotType != wantType {
		t.Errorf("c.GetType() = %d; want %d", gotType, wantType)
	}

}
