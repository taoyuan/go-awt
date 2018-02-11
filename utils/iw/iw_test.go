package iw

import (
	"testing"
	"github.com/smartystreets/goconvey/convey"
	"github.com/prashantv/gostub"
	"io/ioutil"
	"github.com/stretchr/testify/assert"
	"encoding/json"
)

func stubIw(file string) *gostub.Stubs {
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
	convey.Convey("Should scan networks", t, func() {
		cells, err := Scan("")
		if err != nil {
			t.Error(err)
		}
		assert.Equal(t, len(cells), 32)
		actual, _ := json.Marshal(cells)
		//t.Log(string(actual))
		expected, _ := ioutil.ReadFile("test/fixtures/iw-scan-json.txt")
		assert.Equal(t, string(expected), string(actual))
	})
}
