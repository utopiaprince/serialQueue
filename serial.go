package serialQueue

import (
	"bytes"
	"errors"
)

func (sf *SerialFrm) waitSdState(state SerialSigType, ch byte) bool {
	if !sf.register.St.Valid {
		sf.tranState(LD_SIG)
		return false
	}

	if sf.register.St.Data[sf.sdIndex] == ch  {
		sf.sqqueue.Write([]byte{ch})
		sf.sdIndex++
		sf.lastEnterNum++
		if sf.sdIndex >= sf.register.St.len {
			sf.sdIndex = 0
			sf.tranState(LD_SIG)
			return true
		}
	} else {
		if sf.sqqueue.Len() != 0 {
			sf.sqqueue.Reset() //丢弃所有无效数据
		}

		sf.sdIndex = 0
		sf.lastEnterNum = 0
	}

	return true
}

func (sf *SerialFrm) waitLdState(state SerialSigType, ch byte) bool{
	if !sf.register.Ld.Valid {
		sf.tranState(ED_SIG)
		return false
	}

	if sf.lastEnterNum < uint16(sf.register.Ld.Pos) {
		sf.sqqueue.Write([]byte{ch})
		sf.lastEnterNum++
		return true
	}

	sf.ldData[sf.ldIndex] = uint8(ch)
	sf.sqqueue.Write([]byte{ch})
	sf.ldIndex++
	sf.lastEnterNum++
	if sf.ldIndex == SERIAL_LD_LEN_MAX {
		sf.ldVal = uint16(sf.ldData[SERIAL_LD_LEN_MAX-2]) * uint16(256) +
			uint16(sf.ldData[SERIAL_LD_LEN_MAX-1])
		sf.ldVal -= uint16(sf.register.Ld.Pos) + uint16(sf.register.Ld.Len)

		if (sf.ldVal > sf.register.Argu.LenMax) || (sf.ldVal < sf.register.Argu.LenMin) {
			sf.ldIndex = 0
			sf.tranState(SD_SIG)

			if sf.sqqueue.Len() != 0 {
				sf.sqqueue.Reset() //丢弃所有无效数据
			}
			sf.ldIndex = 0
			sf.ldVal = 0
			sf.lastEnterNum = 0
		}

		sf.ldIndex = 0
		sf.tranState(ED_SIG)
	}

	return true
}

func (sf *SerialFrm) endStateHandle() {
	sf.tranState(SD_SIG)

	if sf.fn != nil {
		sf.fn()
	}

	if sf.sqqueue.Len() != 0 {
		sf.sqqueue.Reset() //丢弃所有无效数据
	}
	sf.lastEnterNum = 0
	sf.payloadLen = 0
	sf.ldVal = 0
}

func (sf *SerialFrm) waitEdState(state SerialSigType, ch byte){
	if !sf.register.Ed.Valid {
		if sf.ldVal == 0 {
			sf.endStateHandle()
		}else {
			sf.sqqueue.Write([]byte{ch})
			sf.lastEnterNum++
			sf.payloadLen++
			if sf.payloadLen >= sf.ldVal {
				sf.endStateHandle()
			}
		}
	}else {
		if sf.ldVal == 0 {
			sf.sqqueue.Write([]byte{ch})
			sf.lastEnterNum++
			if sf.register.Ed.DelayEn {
				// TODO: 增加超时回调
			} else {
				// TODO: 判断结束符字符串是否匹配成功，现在是只支持1个字节
				if sf.register.Ed.Data == ch {
					sf.endStateHandle()
				}
			}
		}else {
			sf.sqqueue.Write([]byte{ch})
			sf.lastEnterNum++
			sf.payloadLen++
			if sf.payloadLen >= sf.ldVal {
				if sf.register.Ed.DelayEn {
					// TODO: 增加超时回调
				} else {
					// TODO: 判断结束符字符串是否匹配成功，现在是只支持1个字节
					if sf.register.Ed.Data == ch {
						sf.endStateHandle()
					}
				}
			}
		}
	}
}

func (sf *SerialFrm) Read(n uint16) ([]byte, uint16) {
	p := make([]byte, n)
	m, _ := sf.sqqueue.Read(p)

	return p, uint16(m)
}

func (sf *SerialFrm) FrmLen() uint16{
	sf.Lock()
	defer sf.Unlock()

	return sf.lastEnterNum
}

func (sf *SerialFrm)tranState(state SerialSigType) {
	sf.Lock()
	defer sf.Unlock()

	sf.fsmState = state
}

func (sf *SerialFrm)Write(data []byte) uint16 {
	for _, d := range data {
		if sf.fsmState == SD_SIG {
			if sf.waitSdState(SD_SIG, d) {
				continue
			}
		}

		if sf.fsmState == LD_SIG {
			if sf.waitLdState(LD_SIG, d) {
				continue
			}
		}

		if sf.fsmState == ED_SIG {
			sf.waitEdState(ED_SIG, d)
		}
	}
	return 0
}

func (sf *SerialFrm) AddCallback(cbfunc func()) error{
	if cbfunc == nil {
		return errors.New("callback is null")
	}
	sf.fn = cbfunc

	return nil
}

func New(sReg SerialReg, squeueLen uint16) (*SerialFrm, error) {
	var sFrm SerialFrm

	sFrm.register = &sReg
	sFrm.register.St.len = uint8(len(sFrm.register.St.Data))

	sFrm.ldIndex = 0
	sFrm.ldVal = 0
	sFrm.ldData = make([]byte, sFrm.register.Ld.Len)

	sFrm.payloadLen = 0
	sFrm.edIndex = 0
	sFrm.lastEnterNum = 0

	if sFrm.register.Argu.LenMax == 0 {
		sFrm.register.Argu.LenMax = squeueLen
	}


	sFrm.fsmState = SD_SIG
	sFrm.sqqueue = bytes.NewBuffer(make([]byte, 0, squeueLen))

	return &sFrm, nil
}
