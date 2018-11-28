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
	beego.Info("mySerial: ", mySerial)
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
		readbuf, _ := mySerial.Read(mySerial.FrmLen())
		beego.Info("read buf: ", hex.EncodeToString(readbuf))
	})

	if mySerial == nil {
		t.Errorf("Serial New failed")
	}
	beego.Info("mySerial: ", mySerial)

	originbuf := []byte{0x01, 0x02, 0x0d, 2,3,4, 0x0d}
	beego.Info("write buf: ", hex.EncodeToString(originbuf))
	mySerial.Write(originbuf)
}
