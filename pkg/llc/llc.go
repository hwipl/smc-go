package llc

const (
	// llc/cdc messages are 44 bytes long
	// llc messages are 44 bytes long
	LLCMsgLen = 44
	CDCMsgLen = 44

	// LLC message types
	TypeConfirmLink     = 1
	TypeAddLink         = 2
	TypeAddLinkCont     = 3
	TypeDeleteLink      = 4
	TypeConfirmRKey     = 6
	TypeTestLink        = 7
	TypeConfirmRKeyCont = 8
	TypeDeleteRKey      = 9
	TypeCDC             = 0xFE
)

// ParseLLC parses the LLC message in buffer
func ParseLLC(buffer []byte) Message {
	// llc messages are 44 byte long, treat other lengths as type other
	if len(buffer) != LLCMsgLen {
		return ParseOther(buffer)
	}

	switch buffer[0] {
	case TypeConfirmLink:
		return ParseConfirm(buffer)
	case TypeAddLink:
		return ParseAddLink(buffer)
	case TypeAddLinkCont:
		return ParseAddLinkCont(buffer)
	case TypeDeleteLink:
		return ParseDeleteLink(buffer)
	case TypeConfirmRKey:
		return ParseConfirmRKey(buffer)
	case TypeConfirmRKeyCont:
		return ParseConfirmRKeyCont(buffer)
	case TypeDeleteRKey:
		return parseDeleteRKey(buffer)
	case TypeTestLink:
		return ParseTestLink(buffer)
	case TypeCDC:
		return ParseCDC(buffer)
	default:
		return ParseOther(buffer)
	}
}
