package ifconfig

import (
	"testing"
	"github.com/prashantv/gostub"
	"github.com/smartystreets/goconvey/convey"
	"github.com/stretchr/testify/assert"
	"errors"
)

var IFCONFIG_STATUS_LINUX = `
eth0    Link encap:Ethernet  HWaddr DE:AD:BE:EF:C0:DE
		inet addr:192.168.1.2  Bcast:192.168.1.255  Mask:255.255.255.0
		UP BROADCAST RUNNING MULTICAST  MTU:1500  Metric:1
		RX packets:114919 errors:0 dropped:10 overruns:0 frame:0
		TX packets:117935 errors:0 dropped:0 overruns:0 carrier:0
		collisions:0 txqueuelen:1000
		RX bytes:28178397 (26.8 MiB)  TX bytes:23423409 (22.3 MiB)

lo      Link encap:Local Loopbacks
		inet addr:127.0.0.1  Mask:255.0.0.0
		UP LOOPBACK RUNNING  MTU:65536  Metric:1
		RX packets:0 errors:0 dropped:0 overruns:0 frame:0
		TX packets:0 errors:0 dropped:0 overruns:0 carrier:0
		collisions:0 txqueuelen:0
		RX bytes:0 (0.0 B)  TX bytes:0 (0.0 B)
`

var IFCONFIG_STATUS_INTERFACE_LINUX = `
wlan0   HWaddr DE:AD:BE:EF:C0:DE
		inet6 addr:fe80::21c:c0ff:feae:b5e6/64 Scope:Link
		MTU:1500  Metric:1
		RX packets:0 errors:0 dropped:0 overruns:0 frame:0
		TX packets:0 errors:0 dropped:0 overruns:0 carrier:0
		collisions:0 txqueuelen:1000
		RX bytes:0 (0.0 B)  TX bytes:0 (0.0 B)
`

func TestStatus(t *testing.T) {
	convey.Convey("should get the status for each interface", t, func() {
		stub := gostub.Stub(&Ifconfig, func(args ...string) (string, error) {
			assert.Equal(t, args, []string{"-a"})
			return IFCONFIG_STATUS_LINUX, nil
		})
		defer stub.Reset()

		sts, err := Status("")
		assert.Nil(t, err)
		assert.Equal(t, []*status{
			{
				Iface:          "eth0",
				Link:           "ethernet",
				Address:        "de:ad:be:ef:c0:de",
				Ipv4Address:    "192.168.1.2",
				Ipv4Broadcast:  "192.168.1.255",
				Ipv4SubnetMask: "255.255.255.0",
				Up:             true,
				Broadcast:      true,
				Running:        true,
				Multicast:      true,
			},
			{
				Iface:          "lo",
				Link:           "local",
				Ipv4Address:    "127.0.0.1",
				Ipv4SubnetMask: "255.0.0.0",
				Up:             true,
				Loopback:       true,
				Running:        true,
			},
		}, sts)
	})

	convey.Convey("should get the status for the specified interface", t, func() {
		stub := gostub.Stub(&Ifconfig, func(args ...string) (string, error) {
			assert.Equal(t, args, []string{"wlan0"})
			return IFCONFIG_STATUS_INTERFACE_LINUX, nil
		})
		defer stub.Reset()

		sts, err := Status("wlan0")
		assert.Nil(t, err)
		assert.Equal(t, []*status{
			{
				Iface:       "wlan0",
				Address:     "de:ad:be:ef:c0:de",
				Ipv6Address: "fe80::21c:c0ff:feae:b5e6/64",
			},
		}, sts)
	})

	convey.Convey("should handle error", t, func() {
		err := errors.New("error")
		stub := gostub.Stub(&Ifconfig, func(args ...string) (string, error) {
			return "", err
		})
		defer stub.Reset()

		sts, err2 := Status("")
		assert.Nil(t, sts)
		assert.Equal(t, err, err2)
	})
}
