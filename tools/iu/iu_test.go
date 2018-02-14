package iu

import (
	"testing"
	"github.com/prashantv/gostub"
	"go-awt/infra/osencap"
	"io/ioutil"
	"github.com/smartystreets/goconvey/convey"
	"github.com/stretchr/testify/assert"
	"errors"
	"time"
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
	stub := stubLshw("")
	defer stub.Reset()

	var hw, err = Lshw("")
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, hw.Id, "raspberrypi")
	assert.Equal(t, len(hw.Children), 4)
}

func TestResolveIface(t *testing.T) {
	stub := stubLshw("test/fixtures/lshw-network.txt")
	defer stub.Reset()
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

func TestDown(t *testing.T) {
	convey.Convey("Success", t, func() {
		stub := gostub.Stub(&Ip, func(args ...string) (string, error) {
			assert.Equal(t, args, []string{"link", "set", "wlan0", "down"})
			return "", nil
		})
		defer stub.Reset()

		err := Down("wlan0")
		assert.Nil(t, err)
	})

	convey.Convey("should handle error", t, func() {
		err := errors.New("error")
		stub := gostub.Stub(&Ip, func(args ...string) (string, error) {
			return "", err
		})
		defer stub.Reset()

		err2 := Down("wlan0")
		assert.Equal(t, err, err2)
	})
}

func TestUp(t *testing.T) {
	convey.Convey("Success", t, func() {
		stub := gostub.Stub(&Ip, func(args ...string) (string, error) {
			assert.Equal(t, args, []string{"link", "set", "wlan0", "up"})
			return "", nil
		})
		defer stub.Reset()

		err := Up("wlan0")
		assert.Nil(t, err)
	})

	convey.Convey("should handle error", t, func() {
		err := errors.New("error")
		stub := gostub.Stub(&Ip, func(args ...string) (string, error) {
			return "", err
		})
		defer stub.Reset()

		err2 := Up("wlan0")
		assert.Equal(t, err, err2)
	})
}

func TestReset(t *testing.T) {
	convey.Convey("Success with 0 delay", t, func() {
		ops := []string{"down", "up"}
		i := 0
		stub := gostub.Stub(&Ip, func(args ...string) (string, error) {
			assert.Equal(t, args, []string{"link", "set", "wlan0", ops[i]})
			i = i + 1
			return "", nil
		})
		defer stub.Reset()

		err := Reset("wlan0", 0)
		assert.Nil(t, err)
		assert.Equal(t, 2, i)
	})

	convey.Convey("Success with some delay", t, func() {
		stub := gostub.Stub(&Ip, func(args ...string) (string, error) {
			return "", nil
		})
		defer stub.Reset()

		start := time.Now()
		err := Reset("wlan0", 100)
		end := time.Now()
		delta := end.Sub(start)

		assert.Nil(t, err)
		assert.True(t, delta.Nanoseconds() > 100 * int64(time.Millisecond))
		assert.True(t, delta.Nanoseconds() < 110 * int64(time.Millisecond))
	})

	convey.Convey("should handle error", t, func() {
		err := errors.New("error")
		stub := gostub.Stub(&Ip, func(args ...string) (string, error) {
			return "", err
		})
		defer stub.Reset()

		err2 := Reset("wlan0", 0)
		assert.Equal(t, err, err2)
	})
}