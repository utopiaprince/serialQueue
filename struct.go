package serialQueue

import (
	"bytes"
	"sync"
)

type SerialSigType uint8

const (
	SD_SIG SerialSigType = iota + 0x01
	LD_SIG
	CHAR_SIG
	ED_SIG
)

const SERIAL_LD_LEN_MAX	 = 2

type SerialStart struct {
	len   uint8
	Data  []uint8
	Valid bool
}

type SerialLenDesc struct {
	Len    uint8
	Pos    uint8
	Data   []uint8
	Valid  bool
}

type SerialArgu struct {
	LenMax uint16
	LenMin uint16
}

type SerialEnd struct {
	Len   uint8
	Data  uint8
	Valid bool
	DelayEn  bool
	DelayMs  uint16
}

type SerialReg struct {
	St   SerialStart
	Ld   SerialLenDesc
	Argu SerialArgu
	Ed   SerialEnd
}



type SerialFrm struct {
	sync.RWMutex
	register 	 *SerialReg
	sdIndex      uint8

	ldIndex      uint8
	ldVal        uint16
	ldData      []byte

	payloadLen   uint16
	edIndex      uint8
	lastEnterNum uint16
	locked       uint8
	//sig          uint8
	//char         uint8

	sqqueue      *bytes.Buffer
	fsmState 	 SerialSigType
	fn    		 func()
}



