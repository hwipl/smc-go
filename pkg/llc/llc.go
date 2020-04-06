package llc

const (
	// llc/cdc messages are 44 bytes long
	// llc messages are 44 bytes long
	llcMsgLen = 44
	cdcMsgLen = 44

	// LLC message types
	typeConfirmLink     = 1
	typeAddLink         = 2
	typeAddLinkCont     = 3
	typeDeleteLink      = 4
	typeConfirmRKey     = 6
	typeTestLink        = 7
	typeConfirmRKeyCont = 8
	typeDeleteRKey      = 9
	typeCDC             = 0xFE
)

// ParseLLC parses the LLC message in buffer
func ParseLLC(buffer []byte) Message {
	// llc messages are 44 byte long, treat other lengths as type other
	if len(buffer) != llcMsgLen {
		return parseOther(buffer)
	}

	switch buffer[0] {
	case typeConfirmLink:
		return ParseConfirm(buffer)
	case typeAddLink:
		return ParseAddLink(buffer)
	case typeAddLinkCont:
		return ParseAddLinkCont(buffer)
	case typeDeleteLink:
		return ParseDeleteLink(buffer)
	case typeConfirmRKey:
		return ParseConfirmRKey(buffer)
	case typeConfirmRKeyCont:
		return ParseConfirmRKeyCont(buffer)
	case typeDeleteRKey:
		return parseDeleteRKey(buffer)
	case typeTestLink:
		return ParseTestLink(buffer)
	case typeCDC:
		return ParseCDC(buffer)
	default:
		return parseOther(buffer)
	}
}
