package iw

import (
	"testing"
	"github.com/smartystreets/goconvey/convey"
)

func TestScan(t *testing.T) {
	convey.Convey("Should scan", t, func() {
		cells, err := Scan("")
		if err != nil {
			t.Error(err)
		}
		t.Log(cells)
	})
}