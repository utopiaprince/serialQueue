package serialQueue

import (
	"encoding/hex"
	"github.com/astaxie/beego"
	"testing"
)

func TestSerial_Create(t *testing.T) {

	mySerial, _ := New(SerialReg{}, 1024)

	mySerial.AddCallback(func() {
		t.Logf("mySerial len:%d", mySerial.FrmLen())
	})

	if mySerial == nil {
		t.Errorf("Serial New failed")
	}
}

func TestSerial_Frame_Ed(t *testing.T) {
	mySerial, _ := New(SerialReg{
		Argu: SerialArgu{
			LenMax: 100,
			LenMin: 1,
		},
		Ed: SerialEnd{
			Data:  0x0d,
			Valid: true,
		},
	}, 1024)

	mySerial.AddCallback(func() {
		beego.Info("mySerial len: ", mySerial.FrmLen())
		readBuf, _ := mySerial.Read(mySerial.FrmLen())
		beego.Info("read buf: ", hex.EncodeToString(readBuf))
	})

	if mySerial == nil {
		t.Errorf("Serial New failed")
	}

	originBuf := []byte{
		0x01, 0x02, 0x0d,
		2, 3, 4, 0x0d}
	beego.Info("write buf: ", hex.EncodeToString(originBuf))
	mySerial.Write(originBuf)
}

func TestSerialFrame_St_Ed(t *testing.T) {
	mySerial, _ := New(SerialReg{
		St: SerialStart{
			Valid: true,
			Data:  []uint8{0x01, 0x02},
		},
		Ed: SerialEnd{
			Data:  0x0d,
			Valid: true,
		},
	}, 1024)
	mySerial.AddCallback(func() {
		beego.Info("mySerial len: ", mySerial.FrmLen())
		readBuf, _ := mySerial.Read(mySerial.FrmLen())
		beego.Info("read buf: ", hex.EncodeToString(readBuf))
	})

	if mySerial == nil {
		t.Errorf("Serial New failed")
	}

	originBuf := []byte{
		0x01, 0x02, 0x0d,
		2, 3, 4, 0x0d,
		0x01, 0x02, 0x03, 0x04, 0x0d}
	beego.Info("write buf: ", hex.EncodeToString(originBuf))
	mySerial.Write(originBuf)
}

func TestSerialFrame_St_Ld(t *testing.T) {
	mySerial, _ := New(SerialReg{
		St: SerialStart{
			Valid: true,
			Data:  []uint8{0x01, 0x02},
		},
		Ld: SerialLenDesc{
			Valid: true,
			Pos:   2,
			Len:   2,
		},
		Argu: SerialArgu{
			LenMax: 16,
			LenMin: 2,
		},
	}, 1024)
	mySerial.AddCallback(func() {
		beego.Info("mySerial len: ", mySerial.FrmLen())
		readBuf, _ := mySerial.Read(mySerial.FrmLen())
		beego.Info("read buf: ", hex.EncodeToString(readBuf))
	})

	if mySerial == nil {
		t.Errorf("Serial New failed")
	}
	originBuf := []byte{
		0x01, 0x02, 0x00, 0x06, 3, 4,
		2, 3, 4, 0x0d,
		0x01, 0x02, 0x00, 0x08, 0x0a, 0x0b, 0x0c, 0x0d}
	beego.Info("write buf: ", hex.EncodeToString(originBuf))
	mySerial.Write(originBuf)
}

func TestSerialFrame_St_Ld_Offset(t *testing.T) {
	mySerial, _ := New(SerialReg{
		St: SerialStart{
			Valid: true,
			Data:  []uint8{0x01},
		},
		Ld: SerialLenDesc{
			Valid: true,
			Pos:   2,
			Len:   2,
		},
		Argu: SerialArgu{
			LenMax: 16,
			LenMin: 2,
		},
	}, 1024)
	mySerial.AddCallback(func() {
		beego.Info("mySerial len: ", mySerial.FrmLen())
		readBuf, _ := mySerial.Read(mySerial.FrmLen())
		beego.Info("read buf: ", hex.EncodeToString(readBuf))
	})

	if mySerial == nil {
		t.Errorf("Serial New failed")
	}
	originBuf := []byte{
		0x01, 0x03, 0x00, 0x06, 3, 4,
		2, 3, 4, 0x0d,
		0x01, 0x08, 0x00, 0x08, 0x0a, 0x0b, 0x0c, 0x0d}
	beego.Info("write buf: ", hex.EncodeToString(originBuf))
	mySerial.Write(originBuf)
}
