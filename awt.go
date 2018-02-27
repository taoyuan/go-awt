package awt

import (
	"github.com/taoyuan/go-awt/tools/wpa"
	"strings"
	"github.com/taoyuan/go-awt/tools/iw"
	"github.com/taoyuan/go-awt/tools/iu"
)

func ResolveIface(iface string) (string, error) {
	return iu.ResolveIface(iface)
}

func Status(iface string)(*wpa.Status, error) {
	return wpa.GetStatus(iface)
}

func Mode(iface string) string {
	s, _ := Status(iface)
	state := strings.ToLower(s.WpaState)
	if state == "disconnected" && s.Ip != "" {
		return "ap"
	}
	return "st"
}

func IsAP(iface string) bool {
	return Mode(iface) == "ap"
}

func IsStation(iface string) bool {
	return Mode(iface) == "st"
}

func IsConnected(iface string) bool {
	s, _ := Status(iface)
	state := strings.ToLower(s.WpaState)
	return state == "connected" || state == "completed"
}

func Scan(iface string) iw.Cells {
	cells, _ := iw.Scan(iface)
	return cells
}


