package serialQueue

const (
	SD_SIG uint8 = 0x01
	LD_SIG
	CHAR_SIG
	ED_SIG
)

type SerialStart struct {
	len   uint8
	pos   uint8
	data  []*uint8
	vaild bool
}

type SerialLenDesc struct {
	len    uint8
	pos    uint8
	lenVal uint16
	index  uint8
	data   []*uint8
	vaild  bool
}

type SerialArgu struct {
	lenMax uint16
	lenMin uint16
}

type SerialEnd struct {
	len   uint8
	data  []*uint8
	vaild bool

	delayEn  bool
	delayMs  uint16
	callback func() uint16
}

type SerialReg struct {
	st   SerialStart
	ld   SerialLenDesc
	argu SerialArgu
	sd   SerialEnd
}

type SerialFrm struct {
	sdIndex      uint8
	ld           SerialLenDesc
	payloadLen   uint16
	edIndex      uint8
	lastEnterNum uint16
	locked       uint8
	sig          uint8
	char         uint8
}
