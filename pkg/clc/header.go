package clc

import (
	"encoding/binary"
	"fmt"
)

const (
	// smc type/path
	smcTypeR = 0 // SMC-R only
	smcTypeD = 1 // SMC-D only
	smcTypeB = 3 // SMC-R and SMC-D

	// HeaderLen is the length of the clc header in bytes
	HeaderLen = 8

	// clc message types
	typeProposal = 0x01
	typeAccept   = 0x02
	typeConfirm  = 0x03
	typeDecline  = 0x04
)

// msgType stores the type of a CLC message
type msgType uint8

// String() converts the message type to a string
func (t msgType) String() string {
	switch t {
	case typeProposal:
		return "Proposal"
	case typeAccept:
		return "Accept"
	case typeConfirm:
		return "Confirm"
	case typeDecline:
		return "Decline"
	default:
		return "Unknown"
	}
}

// path stores an SMC path
type path uint8

// String converts the path to a string
func (p path) String() string {
	switch p {
	case smcTypeR:
		return "SMC-R"
	case smcTypeD:
		return "SMC-D"
	case smcTypeB:
		return "SMC-R + SMC-D"
	default:
		return "unknown"
	}
}

// Header stores the common clc message Header
type Header struct {
	// Eyecatcher
	Eyecatcher Eyecatcher

	// type of message: proposal, accept, confirm, decline
	Type msgType

	// total length of message
	Length uint16

	// 1 byte bitfield containing Version, flag, reserved, path:
	Version  uint8 // (4 bits)
	Flag     uint8 // (1 bit)
	reserved byte  // (1 bit)
	Path     path  // (2 bits)
}

// Parse parses the CLC message header in buf
func (h *Header) Parse(buf []byte) {
	// eyecatcher
	copy(h.Eyecatcher[:], buf[:EyecatcherLen])

	// type
	h.Type = msgType(buf[4])

	// length
	h.Length = binary.BigEndian.Uint16(buf[5:7])

	// 1 byte bitfield: version, flag, reserved, path
	bitfield := buf[7]
	h.Version = (bitfield & 0b11110000) >> 4
	h.Flag = (bitfield & 0b00001000) >> 3
	h.reserved = (bitfield & 0b00000100) >> 2
	h.Path = path(bitfield & 0b00000011)
}

// flagString() converts the flag bit in the message according to message type
func (h *Header) flagString() string {
	switch h.Type {
	case typeProposal:
		return fmt.Sprintf("Flag: %d", h.Flag)
	case typeAccept:
		return fmt.Sprintf("First Contact: %d", h.Flag)
	case typeConfirm:
		return fmt.Sprintf("Flag: %d", h.Flag)
	case typeDecline:
		return fmt.Sprintf("Out of Sync: %d", h.Flag)
	default:
		return fmt.Sprintf("Flag: %d", h.Flag)
	}
}

// headerString converts the message header to a string
func (h *Header) String() string {
	flg := h.flagString()
	headerFmt := "%s: Eyecatcher: %s, Type: %d (%s), Length: %d, " +
		"Version: %d, %s, Path: %s"
	return fmt.Sprintf(headerFmt, h.Type, h.Eyecatcher, h.Type, h.Type,
		h.Length, h.Version, flg, h.Path)
}

// Reserved converts the message header fields to a string including reserved
// message fields
func (h *Header) Reserved() string {
	// construct string
	flg := h.flagString()

	headerFmt := "%s: Eyecatcher: %s, Type: %d (%s), Length: %d, " +
		"Version: %d, %s, Reserved: %#x, Path: %s"
	return fmt.Sprintf(headerFmt, h.Type, h.Eyecatcher, h.Type, h.Type,
		h.Length, h.Version, flg, h.reserved, h.Path)
}
