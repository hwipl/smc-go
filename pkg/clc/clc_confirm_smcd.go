package clc

// ConfirmSMCD stores a SMC-D CLC Confirm message
type ConfirmSMCD struct {
	// accept and confirm message have the same message fields
	AcceptSMCD
}
