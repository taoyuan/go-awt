package cap

import (
	"testing"
	"github.com/smartystreets/goconvey/convey"
	"time"
	"github.com/stretchr/testify/assert"
	"github.com/prashantv/gostub"
	"strings"
)

var OUTPUT_STARTED = `WARN: Your adapter does not fully support AP virtual interface, enabling --no-virt
WARN: If AP doesn't work, please read: howto/realtek.md
Config dir: /tmp/create_ap.wlan0.conf.Tg1VJRNW
PID: 3233
Network Manager found, set wlan0 as unmanaged device... DONE
Sharing Internet using method: nat
hostapd command-line interface: hostapd_cli -p /tmp/create_ap.wlan0.conf.Tg1VJRNW/hostapd_ctrl
Configuration file: /tmp/create_ap.wlan0.conf.Tg1VJRNW/hostapd.conf
Using interface wlan0 with hwaddr e8:4e:06:34:ff:db and ssid "MyAccessPoint"
wlan0: interface state UNINITIALIZED->ENABLED
wlan0: AP-ENABLED`

var OUTPUT_DONE = `wlan0: interface state ENABLED->DISABLED
wlan0: AP-DISABLED
nl80211: deinit ifname=wlan0 disabled_11b_rates=0

Doing cleanup.. done`

var OUTPUT_ERROR = `WARN: Your adapter does not fully support AP virtual interface, enabling --no-virt
ERROR: 'eth1' is not an interface`

func TestCreateAP(t *testing.T) {
	convey.Convey("Run success", t, func() {
		stub := gostub.Stub(&CommandCreateAP, "test/create_ap_ok")
		defer stub.Reset()

		c, err := CreateAP(&Options{Iface: "wlan0"})
		assert.Nil(t, err)
		assert.NotNil(t, c)

		c.Start()
		time.Sleep(time.Duration(1) * time.Second)

		output := c.Output()
		expected := strings.Split(OUTPUT_STARTED, "\n")
		assert.Equal(t, expected, output)

		c.Stop()
		time.Sleep(time.Duration(1) * time.Second)

		output = c.Output()
		expected = strings.Split(OUTPUT_DONE, "\n")
		assert.Equal(t, expected, output)

		err = c.Wait()
		assert.Nil(t, err)
	})

	convey.Convey("Run error", t, func() {
		stub := gostub.Stub(&CommandCreateAP, "test/create_ap_error")
		defer stub.Reset()

		c, err := CreateAP(&Options{Iface: "wlan0"})
		assert.Nil(t, err)
		assert.NotNil(t, c)

		c.Start()
		time.Sleep(time.Duration(1) * time.Second)

		output := c.Output()
		expected := strings.Split(OUTPUT_ERROR, "\n")
		assert.Equal(t, expected, output)

		err = c.Wait()
		assert.NotNil(t, err)
	})
}