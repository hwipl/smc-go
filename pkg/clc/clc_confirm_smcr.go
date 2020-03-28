package clc

// confirmSMCR stores a CLC Confirm message
type confirmSMCR struct {
	// accept and confirm messages have the same message fields
	AcceptSMCR
}
