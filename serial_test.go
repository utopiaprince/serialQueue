package serialQueue

import (
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
