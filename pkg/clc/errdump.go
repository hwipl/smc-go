package clc

import (
	"encoding/hex"
	"log"
)

// errDump dumps buffer content in case of an error
func errDump(buf []byte) {
	log.Printf("Message Buffer Hex Dump:\n%s", hex.Dump(buf))
}
