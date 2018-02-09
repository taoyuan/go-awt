package iu

import (
	"testing"
	"github.com/prashantv/gostub"
	"go-awt/infra/osencap"
	"io/ioutil"
	"github.com/magiconair/properties/assert"
	"github.com/smartystreets/goconvey/convey"
)

func stubLshw(file string) *gostub.Stubs  {
	if file == "" {
		file = "test/fixtures/lshw.json"
	}
	return gostub.Stub(&osencap.Exec, func(cmd string, args ...string) (string, error) {
		var raw, err = ioutil.ReadFile(file)
		if err != nil {
			return "", err
		}
		return string(raw), nil

	})
}

func TestLshw(t *testing.T) {
	stubs := stubLshw("")
	defer stubs.Reset()

	var hw, err = Lshw("")
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, hw.Id, "raspberrypi")
	assert.Equal(t, len(hw.Children), 4)
}

func TestResolveIface(t *testing.T) {
	stubs := stubLshw("test/fixtures/lshw-network.txt")
	defer stubs.Reset()
	convey.Convey("Should resolve wlan without calculate", t, func() {
		result, _ := ResolveIface("wlanxyz")
		assert.Equal(t, result, "")
	})

	convey.Convey("Should resolve default iface", t, func() {
		result, _ := ResolveIface("default")
		assert.Equal(t, result, "wlan1")
	})

	convey.Convey("Should resolve onboard iface", t, func() {
		result, _ := ResolveIface("onboard")
		assert.Equal(t, result, "wlan1")
	})

	convey.Convey("Should resolve usb iface", t, func() {
		result, _ := ResolveIface("usb")
		assert.Equal(t, result, "wlan0")
	})

	convey.Convey("Should resolve usb iface", t, func() {
		result, _ := ResolveIface("usb@1:1.2")
		assert.Equal(t, result, "wlan0")
	})
}
