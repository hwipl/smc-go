package clc

import (
	"encoding/hex"
	"log"
	"testing"
)

func TestParseCLCProposalV2SMCB(t *testing.T) {
	// prepare smc-b (r + d) proposal v2 message:
	// Eyecatcher, Type, Length, Version, Pathv2+Path, SenderPeerID
	msgBytes := "e2d4c3d9" + "01" + "00d6" + "2" + "e" +
		"394498039babcdef" +
		// IBGID
		"fe800000000000009a039bfffeabcdef" +
		// IBMAC, IPAreaOffset, SMCDGID
		"98039babcdef" + "0028" + "0123456789abcdef" +
		// ISMv2VCHID, SMCv2Offset, reserved
		"1234" + "0000" + "000000000000000000000000" +
		// reserved
		"00000000000000000000000000000000" +
		// EIDNumber, GIDNumber, reserved3,
		// Release+reserved4+SEIDInd, reserved5, SMCDv2Off
		"01" + "01" + "00" + "01" + "0000" + "00" +
		// SMCDv2Off 1 byte, reserved6 15 bytes
		"40" + "000000000000000000000000000000" +
		// reserved6 16 bytes
		"00000000000000000000000000000000" +
		// reserved6 1 byte, EID 15 bytes
		"00" + "546869734973534d43763245494430" +
		// EID 16 bytes
		"31000000000000000000000000000000" +
		// EID 1 byte, SEID 15 bytes
		"00" + "546869734973534d43763245494430" +
		// SEID 16 bytes
		"32000000000000000000000000000000" +
		// SEID 1 byte, reserved7 15 bytes
		"00" + "000000000000000000000000000000" +
		// reserved7 1 byte, GID 8 bytes, VCHID 2 bytes
		"00" + "abcdef0123456789" + "0123" +
		// Trailer
		"e2d4c3d9"
	msg, err := hex.DecodeString(msgBytes)
	if err != nil {
		log.Fatal(err)
	}

	// parse message
	proposal, proposalLen := NewMessage(msg)
	proposal.Parse(msg)

	// check message length
	if proposalLen != 214 {
		t.Errorf("proposalLen = %d; want %d", proposalLen, 214)
	}

	// check output message without reserved fields
	want := "Proposal: Eyecatcher: SMC-R, Type: 1 (Proposal), " +
		"Length: 214, Version: 2, Pathv2: SMC-R + SMC-D, " +
		"Path: No SMC-R/SMC-D, Peer ID: 14660@98:03:9b:ab:cd:ef, " +
		"SMC-R GID: fe80::9a03:9bff:feab:cdef, " +
		"RoCE MAC: 98:03:9b:ab:cd:ef, IP Area Offset: 40, " +
		"SMC-D GID: 81985529216486895, ISMv2 VCHID: 4660, " +
		"SMCv2 Extension Offset: 0, " +
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
		"Length: 214, Version: 2, Pathv2: SMC-R + SMC-D, " +
		"Path: No SMC-R/SMC-D, Peer ID: 14660@98:03:9b:ab:cd:ef, " +
		"SMC-R GID: fe80::9a03:9bff:feab:cdef, " +
		"RoCE MAC: 98:03:9b:ab:cd:ef, IP Area Offset: 40, " +
		"SMC-D GID: 81985529216486895, ISMv2 VCHID: 4660, " +
		"SMCv2 Extension Offset: 0, " +
		"Reserved: 0x00000000000000000000000000000000000000000000" +
		"000000000000, " +
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

func TestParseCLCProposalV2SMCBIPv4(t *testing.T) {
	// prepare smc-b (r + d) ipv4 proposal v2 message:
	// Eyecatcher, Type, Length, Version, Pathv2+Path, SenderPeerID
	msgBytes := "e2d4c3d9" + "01" + "00de" + "2" + "f" +
		"394498039babcdef" +
		// IBGID
		"fe800000000000009a039bfffeabcdef" +
		// IBMAC, IPAreaOffset, SMCDGID
		"98039babcdef" + "0028" + "0123456789abcdef" +
		// ISMv2VCHID, SMCv2Offset, reserved
		"1234" + "0019" + "000000000000000000000000" +
		// reserved
		"00000000000000000000000000000000" +
		// Prefix, PrefixLen, reserved2, IPv6PrefixesCnt
		"7f000000" + "08" + "0000" + "00" +
		// EIDNumber, GIDNumber, reserved3,
		// Release+reserved4+SEIDInd, reserved5, SMCDv2Off
		"01" + "01" + "00" + "01" + "0000" + "00" +
		// SMCDv2Off 1 byte, reserved6 15 bytes
		"40" + "000000000000000000000000000000" +
		// reserved6 16 bytes
		"00000000000000000000000000000000" +
		// reserved6 1 byte, EID 15 bytes
		"00" + "546869734973534d43763245494430" +
		// EID 16 bytes
		"31000000000000000000000000000000" +
		// EID 1 byte, SEID 15 bytes
		"00" + "546869734973534d43763245494430" +
		// SEID 16 bytes
		"32000000000000000000000000000000" +
		// SEID 1 byte, reserved7 15 bytes
		"00" + "000000000000000000000000000000" +
		// reserved7 1 byte, GID 8 bytes, VCHID 2 bytes
		"00" + "abcdef0123456789" + "0123" +
		// Trailer
		"e2d4c3d9"
	msg, err := hex.DecodeString(msgBytes)
	if err != nil {
		log.Fatal(err)
	}

	// parse message
	proposal, proposalLen := NewMessage(msg)
	proposal.Parse(msg)

	// check message length
	if proposalLen != 222 {
		t.Errorf("proposalLen = %d; want %d", proposalLen, 222)
	}

	// check output message without reserved fields
	want := "Proposal: Eyecatcher: SMC-R, Type: 1 (Proposal), " +
		"Length: 222, Version: 2, Pathv2: SMC-R + SMC-D, " +
		"Path: SMC-R + SMC-D, Peer ID: 14660@98:03:9b:ab:cd:ef, " +
		"SMC-R GID: fe80::9a03:9bff:feab:cdef, " +
		"RoCE MAC: 98:03:9b:ab:cd:ef, IP Area Offset: 40, " +
		"SMC-D GID: 81985529216486895, ISMv2 VCHID: 4660, " +
		"SMCv2 Extension Offset: 25, IPv4 Prefix: 127.0.0.0/8, " +
		"IPv6 Prefix Count: 0, " +
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
		"Length: 222, Version: 2, Pathv2: SMC-R + SMC-D, " +
		"Path: SMC-R + SMC-D, Peer ID: 14660@98:03:9b:ab:cd:ef, " +
		"SMC-R GID: fe80::9a03:9bff:feab:cdef, " +
		"RoCE MAC: 98:03:9b:ab:cd:ef, IP Area Offset: 40, " +
		"SMC-D GID: 81985529216486895, ISMv2 VCHID: 4660, " +
		"SMCv2 Extension Offset: 25, " +
		"Reserved: 0x00000000000000000000000000000000000000000000" +
		"000000000000, IPv4 Prefix: 127.0.0.0/8, Reserved: 0x0000, " +
		"IPv6 Prefix Count: 0, " +
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

func TestParseCLCProposalV2SMCBIPv6(t *testing.T) {
	// prepare smc-b (r + d) ipv6 proposal v2 message:
	// Eyecatcher, Type, Length, Version, Pathv2+Path, SenderPeerID
	ipv6Proposal := "e2d4c3d9" + "01" + "00ef" + "2" + "f" +
		"394498039babcdef" +
		// IBGID
		"fe800000000000009a039bfffeabcdef" +
		// IBMAC, IPAreaOffset, SMCDGID
		"98039babcdef" + "0028" + "0123456789abcdef" +
		// ISMv2VCHID, SMCv2Offset, reserved
		"1234" + "0019" + "000000000000000000000000" +
		// reserved
		"00000000000000000000000000000000" +
		// Prefix, PrefixLen, reserved2, IPv6PrefixesCnt, prefix
		"00000000" + "00" + "0000" + "01" + "0000000000000000" +
		// prefix, prefixLen, EIDNumber, GIDNumber, reserved3,
		// Release+reserved4+SEIDInd, reserved5, SMCDv2Off
		"0000000000000001" + "80" + "01" + "01" + "00" + "01" +
		"0000" + "00" +
		// SMCDv2Off 1 byte, reserved6 15 bytes
		"40" + "000000000000000000000000000000" +
		// reserved6 16 bytes
		"00000000000000000000000000000000" +
		// reserved6 1 byte, EID 15 bytes
		"00" + "546869734973534d43763245494430" +
		// EID 16 bytes
		"31000000000000000000000000000000" +
		// EID 1 byte, SEID 15 bytes
		"00" + "546869734973534d43763245494430" +
		// SEID 16 bytes
		"32000000000000000000000000000000" +
		// SEID 1 byte, reserved7 15 bytes
		"00" + "000000000000000000000000000000" +
		// reserved7 1 byte, GID 8 bytes, VCHID 2 bytes
		"00" + "abcdef0123456789" + "0123" +
		// Trailer
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
