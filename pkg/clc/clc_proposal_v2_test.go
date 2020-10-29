package clc

import (
	"encoding/hex"
	"log"
	"testing"
)

func TestParseCLCProposalV2SMCBIPv6(t *testing.T) {
	// prepare smc-b (r + d) ipv6 proposal v2 message
	ipv6Proposal := "e2d4c3d90100ef2f394498039babcdef" +
		"fe800000000000009a039bfffeabcdef" +
		"98039babcdef00280123456789abcdef" +
		"12340019000000000000000000000000" +
		"00000000000000000000000000000000" +
		"00000000000000010000000000000000" +
		"00000000000000018001010001000000" +
		"40000000000000000000000000000000" +
		"00000000000000000000000000000000" +
		"00546869734973534d43763245494430" +
		"31000000000000000000000000000000" +
		"00546869734973534d43763245494430" +
		"32000000000000000000000000000000" +
		"00000000000000000000000000000000" +
		"00abcdef01234567890123" +
		"e2d4c3d9"
	msg, err := hex.DecodeString(ipv6Proposal)
	if err != nil {
		log.Fatal(err)
	}

	// parse message
	proposal, proposalLen := NewMessage(msg)
	proposal.Parse(msg)

	// check message length
	if proposalLen != 239 {
		t.Errorf("proposalLen = %d; want %d", proposalLen, 239)
	}

	// check output message without reserved fields
	want := "Proposal: Eyecatcher: SMC-R, Type: 1 (Proposal), " +
		"Length: 239, Version: 2, Pathv2: SMC-R + SMC-D, " +
		"Path: SMC-R + SMC-D, Peer ID: 14660@98:03:9b:ab:cd:ef, " +
		"SMC-R GID: fe80::9a03:9bff:feab:cdef, " +
		"RoCE MAC: 98:03:9b:ab:cd:ef, IP Area Offset: 40, " +
		"SMC-D GID: 81985529216486895, ISMv2 VCHID: 4660, " +
		"SMCv2 Extension Offset: 25, IPv4 Prefix: 0.0.0.0/0, " +
		"IPv6 Prefix Count: 1, IPv6 Prefix: ::1/128, " +
		"EID Number: 1, GID Number: 1, Release: 0, " +
		"SEID Indicator: 1, SMC-Dv2 Extension Offset: 64, " +
		"EID Area: [EID 0: ThisIsSMCv2EID01], " +
		"SEID: ThisIsSMCv2EID02, " +
		"GID Area: [GID 0: 12379813738877118345, VCHID 0: 291], " +
		"Trailer: SMC-R"
	got := proposal.String()
	if got != want {
		t.Errorf("proposal.String() = %s; want %s", got, want)
	}

	// check output message with reserved fields
	want = "Proposal: Eyecatcher: SMC-R, Type: 1 (Proposal), " +
		"Length: 239, Version: 2, Pathv2: SMC-R + SMC-D, " +
		"Path: SMC-R + SMC-D, Peer ID: 14660@98:03:9b:ab:cd:ef, " +
		"SMC-R GID: fe80::9a03:9bff:feab:cdef, " +
		"RoCE MAC: 98:03:9b:ab:cd:ef, IP Area Offset: 40, " +
		"SMC-D GID: 81985529216486895, ISMv2 VCHID: 4660, " +
		"SMCv2 Extension Offset: 25, " +
		"Reserved: 0x00000000000000000000000000000000000000000000" +
		"000000000000, IPv4 Prefix: 0.0.0.0/0, Reserved: 0x0000, " +
		"IPv6 Prefix Count: 1, IPv6 Prefix: ::1/128, " +
		"EID Number: 1, GID Number: 1, Reserved: 0x0, Release: 0, " +
		"Reserved: 0x0, SEID Indicator: 1, Reserved: 0x0000, " +
		"SMC-Dv2 Extension Offset: 64, Reserved: 0x0000000000000000" +
		"000000000000000000000000000000000000000000000000, " +
		"EID Area: [EID 0: ThisIsSMCv2EID01], " +
		"SEID: ThisIsSMCv2EID02, " +
		"Reserved: 0x00000000000000000000000000000000, " +
		"GID Area: [GID 0: 12379813738877118345, VCHID 0: 291], " +
		"Trailer: SMC-R"
	got = proposal.Reserved()
	if got != want {
		t.Errorf("proposal.Reserved() = %s; want %s", got, want)
	}
}
