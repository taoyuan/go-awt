package wpa

import (
	"github.com/prashantv/gostub"
	"testing"
	"github.com/smartystreets/goconvey/convey"
	"io/ioutil"
	"github.com/stretchr/testify/assert"
	"strings"
	"errors"
)

var WPA_CLI_STATUS_SILENCE = ""

var WPA_CLI_STATUS_COMPLETED = strings.Join([]string {
	"bssid=2c:f5:d3:02:ea:d9",
	"freq=2412",
	"ssid=Fake-Wifi",
	"id=0",
	"mode=station",
	"pairwise_cipher=CCMP",
	"group_cipher=CCMP",
	"key_mgmt=WPA2-PSK",
	"wpa_state=COMPLETED",
	"ip_address=10.34.141.168",
	"p2p_device_address=e4:28:9c:a8:53:72",
	"address=e4:28:9c:a8:53:72",
	"uuid=e1cda789-8c88-53e8-ffff-31c304580c1e",
}, "\n")

var WPA_CLI_STATUS_4WAY_HANDSHAKE = strings.Join([]string {
	"bssid=2c:f5:d3:02:ea:d9",
	"freq=2412",
	"ssid=Fake-Wifi",
	"id=0",
	"mode=station",
	"pairwise_cipher=CCMP",
	"group_cipher=CCMP",
	"key_mgmt=WPA2-PSK",
	"wpa_state=4WAY_HANDSHAKE",
	"ip_address=10.34.141.168",
	"p2p_device_address=e4:28:9c:a8:53:72",
	"address=e4:28:9c:a8:53:72",
	"uuid=e1cda789-8c88-53e8-ffff-31c304580c1e",
}, "\n")

var WPA_CLI_STATUS_SCANNING = strings.Join([]string {
	"wpa_state=SCANNING",
	"ip_address=10.34.141.168",
	"p2p_device_address=e4:28:9c:a8:53:72",
	"address=e4:28:9c:a8:53:72",
	"uuid=e1cda789-8c88-53e8-ffff-31c304580c1e",
}, "\n")

var WPA_CLI_SCAN_RESULTS = strings.Join([]string {
	"bssid / frequency / signal level / flags / ssid",
	"2c:f5:d3:02:ea:d9	2472	-31	[WPA-PSK-CCMP+TKIP][WPA2-PSK-CCMP+TKIP][ESS]	FakeWifi",
	"2c:f5:d3:02:ea:d9	2472	-31	[WPA-PSK-CCMP+TKIP][WPA2-PSK-CCMP+TKIP][ESS]	FakeWifi2",
}, "\n")

var WPA_CLI_SCAN_NORESULTS = strings.Join([]string {
	"",
}, "\n")

var WPA_CLI_COMMAND_OK = "OK\n"
var WPA_CLI_COMMAND_FAIL = "FAIL\n"
var WPA_CLI_COMMAND_ID = "0\n"


func stubWpa(file string) *gostub.Stubs {
	return gostub.Stub(&Wpa, func(args ...string) (string, error) {
		var raw, err = ioutil.ReadFile(file)
		if err != nil {
			return "", err
		}
		return string(raw), nil
	})
}

func stubWpaWithData(data string) *gostub.Stubs {
	return gostub.Stub(&Wpa, func(args ...string) (string, error) {
		return data, nil
	})
}

func TestStatus(t *testing.T) {
	convey.Convey("status SILENCE", t, func() {
		stub := stubWpaWithData(WPA_CLI_STATUS_SILENCE)
		defer stub.Reset()

		st, err := Status("wlan0")
		if err != nil {
			t.Error(err)
		}
		assert.Nil(t, st)
	})

	convey.Convey("status COMPLETED", t, func() {
		stub := stubWpaWithData(WPA_CLI_STATUS_COMPLETED)
		defer stub.Reset()

		st, err := Status("wlan0")
		if err != nil {
			t.Error(err)
		}
		assert.Equal(t, &status{
			Bssid:            "2c:f5:d3:02:ea:d9",
			Frequency:        2412,
			Mode:             "station",
			KeyMgmt:          "wpa2-psk",
			Ssid:             "Fake-Wifi",
			PairwiseCypher:   "CCMP",
			GroupCipher:      "CCMP",
			P2pDeviceAddress: "e4:28:9c:a8:53:72",
			WpaState:         "COMPLETED",
			Ip:               "10.34.141.168",
			Mac:              "e4:28:9c:a8:53:72",
			Uuid:             "e1cda789-8c88-53e8-ffff-31c304580c1e",
			Id:               0,
		}, st)
	})

	convey.Convey("status 4WAY_HANDSHAKE", t, func() {
		stub := stubWpaWithData(WPA_CLI_STATUS_4WAY_HANDSHAKE)
		defer stub.Reset()

		st, err := Status("wlan0")
		if err != nil {
			t.Error(err)
		}
		assert.Equal(t, &status{
			Bssid:            "2c:f5:d3:02:ea:d9",
			Frequency:        2412,
			Mode:             "station",
			KeyMgmt:          "wpa2-psk",
			Ssid:             "Fake-Wifi",
			PairwiseCypher:   "CCMP",
			GroupCipher:      "CCMP",
			P2pDeviceAddress: "e4:28:9c:a8:53:72",
			WpaState:         "4WAY_HANDSHAKE",
			Ip:               "10.34.141.168",
			Mac:              "e4:28:9c:a8:53:72",
			Uuid:             "e1cda789-8c88-53e8-ffff-31c304580c1e",
			Id:               0,
		}, st)
	})

	convey.Convey("status SCANNING", t, func() {
		stub := stubWpaWithData(WPA_CLI_STATUS_SCANNING)
		defer stub.Reset()

		st, err := Status("wlan0")
		if err != nil {
			t.Error(err)
		}
		assert.Equal(t, &status{
			P2pDeviceAddress: "e4:28:9c:a8:53:72",
			WpaState: "SCANNING",
			Ip: "10.34.141.168",
			Mac: "e4:28:9c:a8:53:72",
			Uuid: "e1cda789-8c88-53e8-ffff-31c304580c1e",
		}, st)
	})

	convey.Convey("should handle error", t, func() {
		err := errors.New("error")
		stub := gostub.Stub(&Wpa, func(args ...string) (string, error) {
			return "", err
		})
		defer stub.Reset()

		st, err2 := Status("wlan0")
		assert.Nil(t, st)
		assert.Equal(t, err, err2)
	})
}

func TestBssid(t *testing.T) {
	convey.Convey("OK SCANNING", t, func() {
		stub := gostub.Stub(&Wpa, func(args ...string) (string, error) {
			assert.Equal(t, args, []string{"-i", "wlan0", "bssid", "Fake-Wifi", "2c:f5:d3:02:ea:89"})
			return WPA_CLI_COMMAND_OK, nil
		})
		defer stub.Reset()

		result, err := Bssid("wlan0", "2c:f5:d3:02:ea:89", "Fake-Wifi")
		if err != nil {
			t.Error(err)
		}
		assert.Equal(t, "OK", result)
	})

	convey.Convey("FAIL SCANNING", t, func() {
		stub := gostub.Stub(&Wpa, func(args ...string) (string, error) {
			assert.Equal(t, args, []string{"-i", "wlan0", "bssid", "2c:f5:d3:02:ea:89", "Fake-Wifi"})
			return WPA_CLI_COMMAND_FAIL, nil
		})
		defer stub.Reset()

		result, err := Bssid("wlan0", "Fake-Wifi", "2c:f5:d3:02:ea:89")
		if err != nil {
			t.Error(err)
		}
		assert.Equal(t, "FAIL", result)
	})

	convey.Convey("should handle error", t, func() {
		err := errors.New("error")
		stub := gostub.Stub(&Wpa, func(args ...string) (string, error) {
			return "", err
		})
		defer stub.Reset()

		result, err2 := Bssid("wlan0", "2c:f5:d3:02:ea:89", "Fake-Wifi")
		assert.Equal(t, "", result)
		assert.Equal(t, err, err2)
	})
}

func TestReassociate(t *testing.T) {
	convey.Convey("OK result", t, func() {
		stub := gostub.Stub(&Wpa, func(args ...string) (string, error) {
			assert.Equal(t, args, []string{"-i", "wlan0", "reassociate"})
			return WPA_CLI_COMMAND_OK, nil
		})
		defer stub.Reset()

		result, err := Reassociate("wlan0")
		assert.Nil(t, err)
		assert.Equal(t, "OK", result)
	})

	convey.Convey("FAIL result", t, func() {
		stub := gostub.Stub(&Wpa, func(args ...string) (string, error) {
			assert.Equal(t, args, []string{"-i", "wlan0", "reassociate"})
			return WPA_CLI_COMMAND_FAIL, nil
		})
		defer stub.Reset()

		result, err := Reassociate("wlan0")
		assert.Nil(t, err)
		assert.Equal(t, "FAIL", result)
	})

	convey.Convey("should handle error", t, func() {
		err := errors.New("error")
		stub := gostub.Stub(&Wpa, func(args ...string) (string, error) {
			return "", err
		})
		defer stub.Reset()

		result, err2 := Reassociate("wlan0")
		assert.Equal(t, "", result)
		assert.Equal(t, err, err2)
	})
}


func TestSet(t *testing.T) {
	convey.Convey("OK result", t, func() {
		stub := gostub.Stub(&Wpa, func(args ...string) (string, error) {
			assert.Equal(t, args, []string{"-i", "wlan0", "set", "ap_scan", "1"})
			return WPA_CLI_COMMAND_OK, nil
		})
		defer stub.Reset()

		result, err := Set("wlan0", "ap_scan", 1)
		assert.Nil(t, err)
		assert.Equal(t, "OK", result)
	})

	convey.Convey("FAIL result", t, func() {
		stub := gostub.Stub(&Wpa, func(args ...string) (string, error) {
			assert.Equal(t, args, []string{"-i", "wlan0", "set", "ap_scan", "1"})
			return WPA_CLI_COMMAND_FAIL, nil
		})
		defer stub.Reset()

		result, err := Set("wlan0", "ap_scan", 1)
		assert.Nil(t, err)
		assert.Equal(t, "FAIL", result)
	})

	convey.Convey("should handle error", t, func() {
		err := errors.New("error")
		stub := gostub.Stub(&Wpa, func(args ...string) (string, error) {
			return "", err
		})
		defer stub.Reset()

		result, err2 := Set("wlan0", "ap_scan", 1)
		assert.Equal(t, "", result)
		assert.Equal(t, err, err2)
	})
}

func TestAddNetwork(t *testing.T) {
	convey.Convey("OK result", t, func() {
		stub := gostub.Stub(&Wpa, func(args ...string) (string, error) {
			assert.Equal(t, args, []string{"-i", "wlan0", "add_network"})
			return WPA_CLI_COMMAND_ID, nil
		})
		defer stub.Reset()

		result, err := AddNetwork("wlan0")
		assert.Nil(t, err)
		assert.Equal(t, "0", result)
	})

	convey.Convey("FAIL result", t, func() {
		stub := gostub.Stub(&Wpa, func(args ...string) (string, error) {
			assert.Equal(t, args, []string{"-i", "wlan0", "add_network"})
			return WPA_CLI_COMMAND_FAIL, nil
		})
		defer stub.Reset()

		result, err := AddNetwork("wlan0")
		assert.Nil(t, err)
		assert.Equal(t, "FAIL", result)
	})

	convey.Convey("should handle error", t, func() {
		err := errors.New("error")
		stub := gostub.Stub(&Wpa, func(args ...string) (string, error) {
			return "", err
		})
		defer stub.Reset()

		result, err2 := AddNetwork("wlan0")
		assert.Equal(t, "", result)
		assert.Equal(t, err, err2)
	})
}

func TestSetNetwork(t *testing.T) {
	convey.Convey("OK result", t, func() {
		stub := gostub.Stub(&Wpa, func(args ...string) (string, error) {
			assert.Equal(t, args, []string{"-i", "wlan0", "set_network", "0", "scan_ssid", "1"})
			return WPA_CLI_COMMAND_OK, nil
		})
		defer stub.Reset()

		result, err := SetNetwork("wlan0", 0, "scan_ssid", "1")
		assert.Nil(t, err)
		assert.Equal(t, "OK", result)
	})

	convey.Convey("FAIL result", t, func() {
		stub := gostub.Stub(&Wpa, func(args ...string) (string, error) {
			assert.Equal(t, args, []string{"-i", "wlan0", "set_network", "0", "fake_variable", "1"})
			return WPA_CLI_COMMAND_FAIL, nil
		})
		defer stub.Reset()

		result, err := SetNetwork("wlan0", 0, "fake_variable", "1")
		assert.Nil(t, err)
		assert.Equal(t, "FAIL", result)
	})

	convey.Convey("should handle error", t, func() {
		err := errors.New("error")
		stub := gostub.Stub(&Wpa, func(args ...string) (string, error) {
			return "", err
		})
		defer stub.Reset()

		result, err2 := SetNetwork("wlan0", 0, "fake_variable", "1")
		assert.Equal(t, "", result)
		assert.Equal(t, err, err2)
	})
}

func TestEnableNetwork(t *testing.T) {
	convey.Convey("OK result", t, func() {
		stub := gostub.Stub(&Wpa, func(args ...string) (string, error) {
			assert.Equal(t, args, []string{"-i", "wlan0", "enable_network", "0"})
			return WPA_CLI_COMMAND_OK, nil
		})
		defer stub.Reset()

		result, err := EnableNetwork("wlan0", 0)
		assert.Nil(t, err)
		assert.Equal(t, "OK", result)
	})

	convey.Convey("FAIL result", t, func() {
		stub := gostub.Stub(&Wpa, func(args ...string) (string, error) {
			assert.Equal(t, args, []string{"-i", "wlan0", "enable_network", "28"})
			return WPA_CLI_COMMAND_FAIL, nil
		})
		defer stub.Reset()

		result, err := EnableNetwork("wlan0", 28)
		assert.Nil(t, err)
		assert.Equal(t, "FAIL", result)
	})

	convey.Convey("should handle error", t, func() {
		err := errors.New("error")
		stub := gostub.Stub(&Wpa, func(args ...string) (string, error) {
			return "", err
		})
		defer stub.Reset()

		result, err2 := EnableNetwork("wlan0", 28)
		assert.Equal(t, "", result)
		assert.Equal(t, err, err2)
	})
}

func TestDisableNetwork(t *testing.T) {
	convey.Convey("OK result", t, func() {
		stub := gostub.Stub(&Wpa, func(args ...string) (string, error) {
			assert.Equal(t, args, []string{"-i", "wlan0", "disable_network", "0"})
			return WPA_CLI_COMMAND_OK, nil
		})
		defer stub.Reset()

		result, err := DisableNetwork("wlan0", 0)
		assert.Nil(t, err)
		assert.Equal(t, "OK", result)
	})

	convey.Convey("FAIL result", t, func() {
		stub := gostub.Stub(&Wpa, func(args ...string) (string, error) {
			assert.Equal(t, args, []string{"-i", "wlan0", "disable_network", "28"})
			return WPA_CLI_COMMAND_FAIL, nil
		})
		defer stub.Reset()

		result, err := DisableNetwork("wlan0", 28)
		assert.Nil(t, err)
		assert.Equal(t, "FAIL", result)
	})

	convey.Convey("should handle error", t, func() {
		err := errors.New("error")
		stub := gostub.Stub(&Wpa, func(args ...string) (string, error) {
			return "", err
		})
		defer stub.Reset()

		result, err2 := DisableNetwork("wlan0", 28)
		assert.Equal(t, "", result)
		assert.Equal(t, err, err2)
	})
}

func TestRemoveNetwork(t *testing.T) {
	convey.Convey("OK result", t, func() {
		stub := gostub.Stub(&Wpa, func(args ...string) (string, error) {
			assert.Equal(t, args, []string{"-i", "wlan0", "remove_network", "0"})
			return WPA_CLI_COMMAND_OK, nil
		})
		defer stub.Reset()

		result, err := RemoveNetwork("wlan0", 0)
		assert.Nil(t, err)
		assert.Equal(t, "OK", result)
	})

	convey.Convey("FAIL result", t, func() {
		stub := gostub.Stub(&Wpa, func(args ...string) (string, error) {
			assert.Equal(t, args, []string{"-i", "wlan0", "remove_network", "28"})
			return WPA_CLI_COMMAND_FAIL, nil
		})
		defer stub.Reset()

		result, err := RemoveNetwork("wlan0", 28)
		assert.Nil(t, err)
		assert.Equal(t, "FAIL", result)
	})

	convey.Convey("should handle error", t, func() {
		err := errors.New("error")
		stub := gostub.Stub(&Wpa, func(args ...string) (string, error) {
			return "", err
		})
		defer stub.Reset()

		result, err2 := RemoveNetwork("wlan0", 28)
		assert.Equal(t, "", result)
		assert.Equal(t, err, err2)
	})
}

func TestSelectNetwork(t *testing.T) {
	convey.Convey("OK result", t, func() {
		stub := gostub.Stub(&Wpa, func(args ...string) (string, error) {
			assert.Equal(t, args, []string{"-i", "wlan0", "select_network", "0"})
			return WPA_CLI_COMMAND_OK, nil
		})
		defer stub.Reset()

		result, err := SelectNetwork("wlan0", 0)
		assert.Nil(t, err)
		assert.Equal(t, "OK", result)
	})

	convey.Convey("FAIL result", t, func() {
		stub := gostub.Stub(&Wpa, func(args ...string) (string, error) {
			assert.Equal(t, args, []string{"-i", "wlan0", "select_network", "28"})
			return WPA_CLI_COMMAND_FAIL, nil
		})
		defer stub.Reset()

		result, err := SelectNetwork("wlan0", 28)
		assert.Nil(t, err)
		assert.Equal(t, "FAIL", result)
	})

	convey.Convey("should handle error", t, func() {
		err := errors.New("error")
		stub := gostub.Stub(&Wpa, func(args ...string) (string, error) {
			return "", err
		})
		defer stub.Reset()

		result, err2 := SelectNetwork("wlan0", 28)
		assert.Equal(t, "", result)
		assert.Equal(t, err, err2)
	})
}

func TestScan(t *testing.T) {
	convey.Convey("OK result", t, func() {
		stub := gostub.Stub(&Wpa, func(args ...string) (string, error) {
			assert.Equal(t, args, []string{"-i", "wlan0", "scan"})
			return WPA_CLI_COMMAND_OK, nil
		})
		defer stub.Reset()

		result, err := Scan("wlan0")
		assert.Nil(t, err)
		assert.Equal(t, "OK", result)
	})

	convey.Convey("should handle error", t, func() {
		err := errors.New("error")
		stub := gostub.Stub(&Wpa, func(args ...string) (string, error) {
			return "", err
		})
		defer stub.Reset()

		result, err2 := Scan("wlan0")
		assert.Equal(t, "", result)
		assert.Equal(t, err, err2)
	})
}

func TestScanResults(t *testing.T) {
	convey.Convey("scan_results NORESULTS", t, func() {
		stub := gostub.Stub(&Wpa, func(args ...string) (string, error) {
			assert.Equal(t, args, []string{"-i", "wlan0", "scan_results"})
			return WPA_CLI_SCAN_NORESULTS, nil
		})
		defer stub.Reset()

		result, err := ScanResults("wlan0")
		assert.Nil(t, err)
		assert.Equal(t, 0, len(result))
	})

	convey.Convey("scan_results COMPLETED", t, func() {
		stub := gostub.Stub(&Wpa, func(args ...string) (string, error) {
			assert.Equal(t, args, []string{"-i", "wlan0", "scan_results"})
			return WPA_CLI_SCAN_RESULTS, nil
		})
		defer stub.Reset()

		result, err := ScanResults("wlan0")
		assert.Nil(t, err)
		assert.Equal(t, []*network{
			{
				Bssid:       "2c:f5:d3:02:ea:d9",
				Frequency:   2472,
				SignalLevel: -31,
				Flags:       "[WPA-PSK-CCMP+TKIP][WPA2-PSK-CCMP+TKIP][ESS]",
				Ssid:        "FakeWifi",
			},
			{
				Bssid:       "2c:f5:d3:02:ea:d9",
				Frequency:   2472,
				SignalLevel: -31,
				Flags:       "[WPA-PSK-CCMP+TKIP][WPA2-PSK-CCMP+TKIP][ESS]",
				Ssid:        "FakeWifi2",
			},
		}, result)
	})

	convey.Convey("should handle error", t, func() {
		err := errors.New("error")
		stub := gostub.Stub(&Wpa, func(args ...string) (string, error) {
			return "", err
		})
		defer stub.Reset()

		result, err2 := ScanResults("wlan0")
		assert.Nil(t, result)
		assert.Equal(t, err, err2)
	})
}


func TestSaveNetwork(t *testing.T) {
	convey.Convey("OK result", t, func() {
		stub := gostub.Stub(&Wpa, func(args ...string) (string, error) {
			assert.Equal(t, args, []string{"-i", "wlan0", "save_config"})
			return WPA_CLI_COMMAND_OK, nil
		})
		defer stub.Reset()

		result, err := SaveConfig("wlan0")
		assert.Nil(t, err)
		assert.Equal(t, "OK", result)
	})

	convey.Convey("FAIL result", t, func() {
		stub := gostub.Stub(&Wpa, func(args ...string) (string, error) {
			assert.Equal(t, args, []string{"-i", "wlan0", "save_config"})
			return WPA_CLI_COMMAND_FAIL, nil
		})
		defer stub.Reset()

		result, err := SaveConfig("wlan0")
		assert.Nil(t, err)
		assert.Equal(t, "FAIL", result)
	})

	convey.Convey("should handle error", t, func() {
		err := errors.New("error")
		stub := gostub.Stub(&Wpa, func(args ...string) (string, error) {
			return "", err
		})
		defer stub.Reset()

		result, err2 := SaveConfig("wlan0")
		assert.Equal(t, "", result)
		assert.Equal(t, err, err2)
	})
}
