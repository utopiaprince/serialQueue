package serialQueue

import (
	"bytes"
	"errors"
)

func (sf *SerialFrm) waitSdState(state SerialSigType, char byte) {

}

func (sf *SerialFrm) waitLdState(state SerialSigType, char byte) {

}

func (sf *SerialFrm) waitEdState(state SerialSigType, char byte) {

}

func (sf *SerialFrm) Read(n uint16) ([]byte, uint16) {


	return []byte{}, 0
}

func (sf *SerialFrm) FrmLen() uint16{
	sf.Lock()
	defer sf.Unlock()

	return sf.lastEnterNum
}

func (sf *SerialFrm)Write(data []byte) uint16 {
	for _, d := range data {
		if sf.fsm.state == SD_SIG {
			sf.waitSdState(SD_SIG, d)
		}
		if sf.fsm.state == LD_SIG {
			sf.waitLdState(LD_SIG, d)
		}
		if sf.fsm.state == ED_SIG {
			sf.waitEdState(ED_SIG, d)
		}
	}
	return 0
}

func New(sReg SerialReg, squeueLen uint16, cbfunc func(interface{})) (*SerialFrm, error) {
	var sFrm SerialFrm
	if cbfunc == nil {
		return nil, errors.New("callback is null")
	}
	sFrm.register = &sReg
	sFrm.fsm.state = SD_SIG
	sFrm.fsm.fn = cbfunc
	sFrm.sqqueue = bytes.NewBuffer(make([]byte, 0, squeueLen))

	return &sFrm, nil
}