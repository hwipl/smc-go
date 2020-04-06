package llc

import (
	"encoding/hex"
	"log"
	"strings"
	"testing"
)

func TestCDC(t *testing.T) {
	var want, got string

	// create bytes of a message
	msgHex := "00000000  fe 2c 00 01 00 00 00 03  " +
		"00 00 00 00 00 00 00 0d  |.,..............|\n" +
		"00000010  00 00 00 00 00 00 00 00  " +
		"00 00 00 00 00 00 00 00  |................|\n" +
		"00000020  00 00 00 00 00 00 00 00  " +
		"00 00 00 00              |............|\n"
	msg := "fe 2c 00 01 00 00 00 03  00 00 00 00 00 00 00 0d" +
		"00 00 00 00 00 00 00 00  00 00 00 00 00 00 00 00" +
		"00 00 00 00 00 00 00 00  00 00 00 00"
	bytes, err := hex.DecodeString(strings.Join(strings.Fields(msg), ""))
	if err != nil {
		log.Fatal(err)
	}

	// parse message
	cdc := ParseCDC(bytes)

	// test String()
	want = "CDC: Type: 254, Length 44, Sequence Number: 1, " +
		"Alert Token: 3, Producer Wrap: 0, Producer Cursor: 13, " +
		"Consumer Wrap: 0, Consumer Cursor: 0, " +
		"Writer Blocked: false, Urgent Data Pending: false, " +
		"Urgent Data Present: false, " +
		"Request for Consumer Cursor Update: false, " +
		"Failover Validation: false, Sending Done: false, " +
		"Peer Connection Closed: false, Abnormal Close: false\n"
	got = cdc.String()
	if got != want {
		t.Errorf("cdc.String() = %s; want %s", got, want)
	}

	// test Reserved()
	want = "CDC: Type: 254, Length 44, Sequence Number: 1, " +
		"Alert Token: 3, Reserved: 0x0000, Producer Wrap: 0, " +
		"Producer Cursor: 13, Reserved: 0x0000, Consumer Wrap: 0, " +
		"Consumer Cursor: 0, Writer Blocked: false, " +
		"Urgent Data Pending: false, Urgent Data Present: false, " +
		"Request for Consumer Cursor Update: false, " +
		"Failover Validation: false, Reserved: 0x0, " +
		"Sending Done: false, Peer Connection Closed: false, " +
		"Abnormal Close: false, " +
		"Reserved: 0x00000000000000000000000000000000000000\n"
	got = cdc.Reserved()
	if got != want {
		t.Errorf("cdc.Reserved() = %s; want %s", got, want)
	}

	// test Hex()
	want = msgHex
	got = cdc.Hex()
	if got != want {
		t.Errorf("cdc.Hex() = %s; want %s", got, want)
	}

	// test GetType()
	wantType := TypeCDC
	gotType := cdc.GetType()
	if gotType != wantType {
		t.Errorf("cdc.GetType() = %d; want %d", gotType, wantType)
	}

}
