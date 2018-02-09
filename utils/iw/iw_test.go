package iw

import (
	"testing"
	"github.com/smartystreets/goconvey/convey"
	"github.com/prashantv/gostub"
	"io/ioutil"
)

func stubIw(file string) *gostub.Stubs  {
	return gostub.Stub(&Iw, func(args ...string) (string, error) {
		var raw, err = ioutil.ReadFile(file)
		if err != nil {
			return "", err
		}
		return string(raw), nil

	})
}

func TestScan(t *testing.T) {
	stub := stubIw("test/fixtures/iw-scan.txt")
	defer stub.Reset()
	convey.Convey("Should scan", t, func() {
		cells, err := Scan("")
		if err != nil {
			t.Error(err)
		}
		t.Log(cells)
	})
}