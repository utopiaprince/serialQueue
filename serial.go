package serialQueue

import (
	"bytes"
	"errors"
)

func (sf *SerialFrm) waitSdState(state SerialSigType, ch byte) {
	if !sf.register.St.valid {
		sf.tranState(LD_SIG)
		return
	}

	if *sf.register.St.data[sf.sdIndex] == ch  {
		sf.sqqueue.Write([]byte{ch})
		sf.sdIndex++
		sf.lastEnterNum++
		if sf.sdIndex == sf.register.St.len {
			sf.sdIndex = 0
			sf.tranState(LD_SIG)
		}
	} else {
		if sf.sqqueue.Len() != 0 {
			sf.sqqueue.Reset() //丢弃所有无效数据
		}

		sf.sdIndex = 0
		sf.lastEnterNum = 0
	}
}

func (sf *SerialFrm) waitLdState(state SerialSigType, ch byte) {
	if !sf.register.Ld.valid {
		sf.tranState(ED_SIG)
		return
	}

	*sf.ld.data[sf.ld.index] = uint8(ch)
	sf.sqqueue.Write([]byte{ch})
	sf.ld.index++
	sf.lastEnterNum++
	if sf.ld.index == SERIAL_LD_LEN_MAX {
		sf.ld.lenVal = uint16(*sf.ld.data[SERIAL_LD_LEN_MAX-2]) * uint16(256) +
			uint16(*sf.ld.data[SERIAL_LD_LEN_MAX-1]);

		if (sf.ld.lenVal > sf.register.Argu.lenMax) || (sf.ld.lenVal < sf.register.Argu.lenMin) {
			sf.ld.index = 0
			sf.tranState(SD_SIG)

			if sf.sqqueue.Len() != 0 {
				sf.sqqueue.Reset() //丢弃所有无效数据
			}
			sf.ld.index = 0
			sf.ld.lenVal = 0
			sf.lastEnterNum = 0
		}

		sf.ld.index = 0
		sf.tranState(ED_SIG)
	}
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
	sf.ld.lenVal = 0
}

func (sf *SerialFrm) waitEdState(state SerialSigType, ch byte) {
	if !sf.register.Ed.valid {
		if sf.ld.lenVal == 0 {
			sf.endStateHandle()
		}else {
			sf.sqqueue.Write([]byte{ch})
			sf.lastEnterNum++
			sf.payloadLen++
			if sf.payloadLen >= sf.ld.lenVal {
				sf.endStateHandle()
			}
		}
	}else {
		if sf.ld.lenVal == 0 {
			sf.sqqueue.Write([]byte{ch})
			sf.lastEnterNum++
			if sf.register.Ed.delayEn {
				// TODO: 增加超时回调
			} else {
				// TODO: 判断结束符字符串是否匹配成功，现在是只支持1个字节
				if sf.register.Ed.data == ch {
					sf.endStateHandle()
				}
			}
		}else {
			sf.sqqueue.Write([]byte{ch})
			sf.lastEnterNum++
			sf.payloadLen++
			if sf.payloadLen >= sf.ld.lenVal {
				if sf.register.Ed.delayEn {
					// TODO: 增加超时回调
				} else {
					// TODO: 判断结束符字符串是否匹配成功，现在是只支持1个字节
					if sf.register.Ed.data == ch {
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
			sf.waitSdState(SD_SIG, d)
		}
		if sf.fsmState == LD_SIG {
			sf.waitLdState(LD_SIG, d)
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
	sFrm.fsmState = SD_SIG
	sFrm.sqqueue = bytes.NewBuffer(make([]byte, 0, squeueLen))

	return &sFrm, nil
}
