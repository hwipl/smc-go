package llc

import (
	"encoding/hex"
	"log"
	"strings"
	"testing"
)

func TestOther(t *testing.T) {
	var want, got string

	// create bytes of a message
	msgHex := "00000000  0a 00 00 01              " +
		"                         |....|\n"
	msg := "0a 00 00 01"
	bytes, err := hex.DecodeString(strings.Join(strings.Fields(msg), ""))
	if err != nil {
		log.Fatal(err)
	}

	// parse message
	o := ParseOther(bytes)

	// test String()
	want = "Other Payload\n"
	got = o.String()
	if got != want {
		t.Errorf("o.String() = %s; want %s", got, want)
	}

	// test Reserved()
	want = "Other Payload\n"
	got = o.Reserved()
	if got != want {
		t.Errorf("o.Reserved() = %s; want %s", got, want)
	}

	// test Hex()
	want = msgHex
	got = o.Hex()
	if got != want {
		t.Errorf("o.Hex() = %s; want %s", got, want)
	}

	// test GetType()
	wantType := TypeOther
	gotType := o.GetType()
	if gotType != wantType {
		t.Errorf("o.GetType() = %d; want %d", gotType, wantType)
	}

}
