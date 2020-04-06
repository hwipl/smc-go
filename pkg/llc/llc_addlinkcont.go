package llc

import (
	"encoding/binary"
	"fmt"
)

// RKeyPair stores a RKey/RToken pair
type RKeyPair struct {
	referenceRKey uint32
	newRKey       uint32
	newVAddr      uint64
}

// Parse fills the rkeyPair fields from the buffer
func (r *RKeyPair) Parse(buffer []byte) {
	// RKey/RToken pairs are 16 bytes and consist of:
	// * Reference Key (4 bytes)
	// * New RKey (4 bytes)
	// * New Virtual Address (8 bytes)
	r.referenceRKey = binary.BigEndian.Uint32(buffer[0:4])
	r.newRKey = binary.BigEndian.Uint32(buffer[4:8])
	r.newVAddr = binary.BigEndian.Uint64(buffer[8:16])
}

// String converts the rkeyPair to a string
func (r *RKeyPair) String() string {
	rFmt := "[Reference RKey: %d, New RKey: %d, New Virtual Address: %#x]"
	return fmt.Sprintf(rFmt, r.referenceRKey, r.newRKey, r.newVAddr)
}

// AddLinkCont stores a LLC add link continuation message
type AddLinkCont struct {
	BaseMsg
	res1       byte
	Reply      bool
	res2       byte
	Link       uint8
	NumRTokens uint8
	res3       [2]byte
	RKeyPairs  [2]RKeyPair
	res4       [4]byte
}

// Parse fills the addLinkCont fields from the LLC add link continuation
// message in buffer
func (a *AddLinkCont) Parse(buffer []byte) {
	// init base message fields
	a.SetBaseMsg(buffer)
	buffer = buffer[2:]

	// Reserved 1 byte
	a.res1 = buffer[0]
	buffer = buffer[1:]

	// Reply is first bit in this byte
	a.Reply = (buffer[0] & 0b10000000) > 0

	// Remainder of this byte is reserved
	a.res2 = buffer[0] & 0b01111111
	buffer = buffer[1:]

	// Link is 1 byte
	a.Link = buffer[0]
	buffer = buffer[1:]

	// Number of RTokens is 1 byte
	a.NumRTokens = buffer[0]
	buffer = buffer[1:]

	// Reserved 2 bytes
	copy(a.res3[:], buffer[0:2])
	buffer = buffer[2:]

	// RKey/RToken pairs are each 16 bytes
	// parse
	// * first RKey/RToken pair
	// * second RKey/RToken pair (can be all zero)
	for i := range a.RKeyPairs {
		a.RKeyPairs[i].Parse(buffer)
		buffer = buffer[16:]
	}

	// Rest of message is reserved
	copy(a.res4[:], buffer[:])
}

// String converts the add link continuation message to a string
func (a *AddLinkCont) String() string {
	var pairs string

	// convert RKey pairs
	for i := range a.RKeyPairs {
		pairs = fmt.Sprintf(", RKey Pair %d: %s", i+1, &a.RKeyPairs[i])
	}

	aFmt := "LLC Add Link Continuation: Type: %d, Length: %d, " +
		"Reply: %t, Link: %d, Number of RTokens: %d%s\n"
	return fmt.Sprintf(aFmt, a.Type, a.Length, a.Reply, a.Link,
		a.NumRTokens, pairs)
}

// Reserved converts the add link continuation message to a string including
// reserved fields
func (a *AddLinkCont) Reserved() string {
	var pairs string

	// convert RKey pairs
	for i := range a.RKeyPairs {
		pairs = fmt.Sprintf("RKey Pair %d: %s, ", i+1, &a.RKeyPairs[i])
	}

	aFmt := "LLC Add Link Continuation: Type: %d, Length: %d, " +
		"Reserved: %#x, Reply: %t, Reserved: %#x, Link: %d, " +
		"Number of RTokens: %d, Reserved: %#x, %sReserved: %#x\n"
	return fmt.Sprintf(aFmt, a.Type, a.Length, a.res1, a.Reply, a.res2,
		a.Link, a.NumRTokens, a.res3, pairs, a.res4)
}

// ParseAddLinkCont parses the LLC add link continuation message in buffer
func ParseAddLinkCont(buffer []byte) *AddLinkCont {
	var addCont AddLinkCont
	addCont.Parse(buffer)
	return &addCont
}
