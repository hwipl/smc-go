package roce

import (
	"encoding/hex"
	"log"
	"strings"
	"testing"
)

func TestGRH(t *testing.T) {
	var want, got string

	// create bytes of a message
	msg := "60 20 00 00 00 3c 1b 01  fe 80 00 00 00 00 00 00" +
		"9a 03 9b ff fe ab cd ef  fe 80 00 00 00 00 00 00" +
		"9a 03 9b ff fe ab cd ef"
	bytes, err := hex.DecodeString(strings.Join(strings.Fields(msg), ""))
	if err != nil {
		log.Fatal(err)
	}

	// parse message
	grh := ParseGRH(bytes)

	// test String()
	want = "GRH: Version: 6, Traffic Class: 2, Flow Label: 0, " +
		"Length: 60, Next Header: BTH, Hop Limit: 1, " +
		"Source: fe80::9a03:9bff:feab:cdef, " +
		"Destination: fe80::9a03:9bff:feab:cdef\n"
	got = grh.String()
	if got != want {
		t.Errorf("grh.String() = %s; want %s", got, want)
	}

	// test Reserved(), should produce same output because there are no
	// reserved values printed
	got = grh.Reserved()
	if got != want {
		t.Errorf("grh.Reserved() = %s; want %s", got, want)
	}

	// test Hex()
	want = "00000000  60 20 00 00 00 3c 1b 01  " +
		"fe 80 00 00 00 00 00 00  |` ...<..........|\n" +
		"00000010  9a 03 9b ff fe ab cd ef  " +
		"fe 80 00 00 00 00 00 00  |................|\n" +
		"00000020  9a 03 9b ff fe ab cd ef  " +
		"                         |........|\n"
	got = grh.Hex()
	if got != want {
		t.Errorf("grh.Hex() = %s; want %s", got, want)
	}

	// test GetType()
	wantType := typeGRH
	gotType := grh.GetType()
	if gotType != wantType {
		t.Errorf("grh.GetType() = %d; want %d", gotType, wantType)
	}

}
