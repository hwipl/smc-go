package roce

import (
	"encoding/hex"
	"log"
	"strings"
	"testing"
)

func TestBTH(t *testing.T) {
	var want, got string

	// create bytes of a message
	msgHex := "00000000  04 c0 ff ff 00 00 01 07  " +
		"80 7b 28 94              |.........{(.|\n"
	msg := "04 c0 ff ff 00 00 01 07  80 7b 28 94"
	bytes, err := hex.DecodeString(strings.Join(strings.Fields(msg), ""))
	if err != nil {
		log.Fatal(err)
	}

	// parse message
	bth := ParseBTH(bytes)

	// test String()
	want = "BTH: OpCode: 0b100 (RC SEND Only), SE: true, M: true, " +
		"Pad: 0, TVer: 0, PKey: 65535, FECN: false, BECN: false, " +
		"DestQP: 263, A: true, PSN: 8071316\n"
	got = bth.String()
	if got != want {
		t.Errorf("bth.String() = %s; want %s", got, want)
	}

	// test Reserved()
	want = "BTH: OpCode: 0b100 (RC SEND Only), SE: true, M: true, " +
		"Pad: 0, TVer: 0, PKey: 65535, FECN: false, BECN: false, " +
		"Res: 0x0, DestQP: 263, A: true, Res: 0x0, PSN: 8071316\n"
	got = bth.Reserved()
	if got != want {
		t.Errorf("bth.Reserved() = %s; want %s", got, want)
	}

	// test Hex()
	want = msgHex
	got = bth.Hex()
	if got != want {
		t.Errorf("bth.Hex() = %s; want %s", got, want)
	}

	// test GetType()
	wantType := typeBTH
	gotType := bth.GetType()
	if gotType != wantType {
		t.Errorf("bth.GetType() = %d; want %d", gotType, wantType)
	}

}
