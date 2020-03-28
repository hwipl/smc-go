package clc

// confirmSMCD stores a SMC-D CLC Confirm message
type confirmSMCD struct {
	// accept and confirm message have the same message fields
	AcceptSMCD
}
