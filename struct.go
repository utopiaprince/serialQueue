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
	data  []uint8
	valid bool
}

type SerialLenDesc struct {
	len    uint8
	pos    uint8
	data   []uint8
	valid  bool
}

type SerialArgu struct {
	lenMax uint16
	lenMin uint16
}

type SerialEnd struct {
	len   uint8
	data  uint8
	valid bool
	delayEn  bool
	delayMs  uint16
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



