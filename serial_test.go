package serialQueue

import (
	"encoding/hex"
	"testing"
 	"github.com/astaxie/beego"
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

func TestSerial_Frame_Ed(t *testing.T){
	mySerial, _ := New(SerialReg{
		Argu: SerialArgu{
			lenMax: 100,
			lenMin: 1,
		},
		Ed: SerialEnd{
			data: 0x0d,
			valid: true,
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
		2,3,4, 0x0d}
	beego.Info("write buf: ", hex.EncodeToString(originBuf))
	mySerial.Write(originBuf)
}

func TestSerialFrame_St_Ed(t *testing.T) {
	mySerial, _ := New(SerialReg{
		St: SerialStart{
			valid: true,
			data: []uint8{0x01,0x02},
		},
		Ed: SerialEnd{
			data: 0x0d,
			valid: true,
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
		0x01,0x02,0x0d,
		2,3,4,0x0d,
		0x01,0x02,0x03,0x04,0x0d}
	beego.Info("write buf: ", hex.EncodeToString(originBuf))
	mySerial.Write(originBuf)
}

func TestSerialFrame_St_Ld(t *testing.T) {
	mySerial, _ := New(SerialReg{
		St: SerialStart{
			valid: true,
			data: []uint8{0x01,0x02},
		},
		Ld: SerialLenDesc{
			valid: true,
			pos: 2,
			len: 2,
		},
		Argu: SerialArgu{
			lenMax: 16,
			lenMin: 2,
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
		0x01,0x02,0x00,0x02,3,4,
		2,3,4,0x0d,
		0x01,0x02,0x00,0x04,0x0a,0x0b,0x0c,0x0d}
	beego.Info("write buf: ", hex.EncodeToString(originBuf))
	mySerial.Write(originBuf)
}