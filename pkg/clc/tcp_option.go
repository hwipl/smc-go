package clc

import (
	"bytes"

	"github.com/google/gopacket/layers"
)

var (
	// SMCOption is the content of the tcp experimental option
	SMCOption = SMCREyecatcher
)

// CheckSMCOption checks if SMC option is set in TCP header
func CheckSMCOption(tcp *layers.TCP) bool {
	for _, opt := range tcp.Options {
		if opt.OptionType == 254 &&
			opt.OptionLength == 6 &&
			bytes.Compare(opt.OptionData, SMCOption) == 0 {
			return true
		}
	}

	return false
}
