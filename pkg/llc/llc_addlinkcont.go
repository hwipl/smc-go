package llc

import (
	"encoding/binary"
	"fmt"
)

// rkeyPair stores a RKey/RToken pair
type rkeyPair struct {
	referenceRKey uint32
	newRKey       uint32
	newVAddr      uint64
}

// parse fills the rkeyPair fields from the buffer
func (r *rkeyPair) parse(buffer []byte) {
	// RKey/RToken pairs are 16 bytes and consist of:
	// * Reference Key (4 bytes)
	// * New RKey (4 bytes)
	// * New Virtual Address (8 bytes)
	r.referenceRKey = binary.BigEndian.Uint32(buffer[0:4])
	r.newRKey = binary.BigEndian.Uint32(buffer[4:8])
	r.newVAddr = binary.BigEndian.Uint64(buffer[8:16])
}

// String converts the rkeyPair to a string
func (r *rkeyPair) String() string {
	rFmt := "[Reference RKey: %d, New RKey: %d, New Virtual Address: %#x]"
	return fmt.Sprintf(rFmt, r.referenceRKey, r.newRKey, r.newVAddr)
}

// addLinkCont stores a LLC add link continuation message
type addLinkCont struct {
	baseMsg
	res1       byte
	reply      bool
	res2       byte
	link       uint8
	numRTokens uint8
	res3       [2]byte
	rkeyPairs  [2]rkeyPair
	res4       [4]byte
}

// Parse fills the addLinkCont fields from the LLC add link continuation
// message in buffer
func (a *addLinkCont) Parse(buffer []byte) {
	// init base message fields
	a.setBaseMsg(buffer)
	buffer = buffer[2:]

	// Reserved 1 byte
	a.res1 = buffer[0]
	buffer = buffer[1:]

	// Reply is first bit in this byte
	a.reply = (buffer[0] & 0b10000000) > 0

	// Remainder of this byte is reserved
	a.res2 = buffer[0] & 0b01111111
	buffer = buffer[1:]

	// Link is 1 byte
	a.link = buffer[0]
	buffer = buffer[1:]

	// Number of RTokens is 1 byte
	a.numRTokens = buffer[0]
	buffer = buffer[1:]

	// Reserved 2 bytes
	copy(a.res3[:], buffer[0:2])
	buffer = buffer[2:]

	// RKey/RToken pairs are each 16 bytes
	// parse
	// * first RKey/RToken pair
	// * second RKey/RToken pair (can be all zero)
	for i := range a.rkeyPairs {
		a.rkeyPairs[i].parse(buffer)
		buffer = buffer[16:]
	}

	// Rest of message is reserved
	copy(a.res4[:], buffer[:])
}

// String converts the add link continuation message to a string
func (a *addLinkCont) String() string {
	var pairs string

	// convert RKey pairs
	for i := range a.rkeyPairs {
		pairs = fmt.Sprintf(", RKey Pair %d: %s", i+1, &a.rkeyPairs[i])
	}

	aFmt := "LLC Add Link Continuation: Type: %d, Length: %d, " +
		"Reply: %t, Link: %d, Number of RTokens: %d%s\n"
	return fmt.Sprintf(aFmt, a.typ, a.length, a.reply, a.link,
		a.numRTokens, pairs)
}

// Reserved converts the add link continuation message to a string including
// reserved fields
func (a *addLinkCont) Reserved() string {
	var pairs string

	// convert RKey pairs
	for i := range a.rkeyPairs {
		pairs = fmt.Sprintf("RKey Pair %d: %s, ", i+1, &a.rkeyPairs[i])
	}

	aFmt := "LLC Add Link Continuation: Type: %d, Length: %d, " +
		"Reserved: %#x, Reply: %t, Reserved: %#x, Link: %d, " +
		"Number of RTokens: %d, Reserved: %#x, %sReserved: %#x\n"
	return fmt.Sprintf(aFmt, a.typ, a.length, a.res1, a.reply, a.res2,
		a.link, a.numRTokens, a.res3, pairs, a.res4)
}

// parseAddLinkCont parses the LLC add link continuation message in buffer
func parseAddLinkCont(buffer []byte) *addLinkCont {
	var addCont addLinkCont
	addCont.Parse(buffer)
	return &addCont
}
