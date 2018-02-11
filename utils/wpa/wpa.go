package wpa

import (
	"go-awt/infra/osencap"
	"strings"
	"go-awt/infra/re"
	"strconv"
	"github.com/thoas/go-funk"
	"go-awt/infra/str"
)

type status struct {
	Id               int
	Ssid             string
	Bssid            string
	Frequency        int
	Mode             string
	KeyMgmt          string
	PairwiseCypher   string
	GroupCipher      string
	P2pDeviceAddress string
	WpaState         string
	Ip               string
	Mac              string
	Uuid             string
}

type network struct {
	Bssid       string
	Frequency   int
	SignalLevel int
	Flags       string
	Ssid        string
}

var Wpa = func(args ...string) (string, error) {
	return osencap.Exec("wpa_cli", args...)
}

func parseStatusBlock(block string) *status {
	v := false
	st := status{}
	var matches []string

	if matches = re.Match(`bssid=([A-Fa-f0-9:]{17})`, block); len(matches) > 1 {
		st.Bssid = strings.ToLower(matches[1])
		v = true
	}

	if matches = re.Match(`freq=([0-9]+)`, block); len(matches) > 1 {
		st.Frequency, _ = strconv.Atoi(matches[1])
		v = true
	}

	if matches = re.Match(`mode=([^\s]+)`, block); len(matches) > 1 {
		st.Mode = matches[1]
		v = true
	}

	if matches = re.Match(`key_mgmt=([^\s]+)`, block); len(matches) > 1 {
		st.KeyMgmt = strings.ToLower(matches[1])
		v = true
	}

	if matches = re.Match(`[^b]ssid=([^\n]+)`, block); len(matches) > 1 {
		st.Ssid = matches[1]
		v = true
	}

	if matches = re.Match(`[^b]pairwise_cipher=([^\n]+)`, block); len(matches) > 1 {
		st.PairwiseCypher = matches[1]
		v = true
	}

	if matches = re.Match(`[^b]group_cipher=([^\n]+)`, block); len(matches) > 1 {
		st.GroupCipher = matches[1]
		v = true
	}

	if matches = re.Match(`p2p_device_address=([A-Fa-f0-9:]{17})`, block); len(matches) > 1 {
		st.P2pDeviceAddress = matches[1]
		v = true
	}

	if matches = re.Match(`wpa_state=([^\s]+)`, block); len(matches) > 1 {
		st.WpaState = matches[1]
		v = true
	}

	if matches = re.Match(`ip_address=([^\n]+)`, block); len(matches) > 1 {
		st.Ip = matches[1]
		v = true
	}

	if matches = re.Match(`[^_]address=([A-Fa-f0-9:]{17})`, block); len(matches) > 1 {
		st.Mac = matches[1]
		v = true
	}

	if matches = re.Match(`uuid=([^\n]+)`, block); len(matches) > 1 {
		st.Uuid = matches[1]
		v = true
	}

	if matches = re.Match(`[^s]id=([0-9]+)`, block); len(matches) > 1 {
		st.Id, _ = strconv.Atoi(matches[1])
		v = true
	}

	if v {
		return &st
	}
	return nil
}

func parseScanResultsBlock(block string) []*network {
	var answer []*network
	var matches []string

	lines := strings.Split(block, "\n")
	lines = funk.Map(lines, func (line string) string { return line + "\n"} ).([]string)
	for _, entry := range lines {
		v := false
		n := network{}

		if matches = re.Match(`([A-Fa-f0-9:]{17})\t`, entry); len(matches) > 1 {
			n.Bssid = strings.ToLower(matches[1])
			v = true
		}

		if matches = re.Match(`\t([\d]+)\t+`, entry); len(matches) > 1 {
			n.Frequency, _ = strconv.Atoi(matches[1])
			v = true
		}

		if matches = re.Match(`([-][0-9]+)\t`, entry); len(matches) > 1 {
			n.SignalLevel, _ = strconv.Atoi(matches[1])
			v = true
		}

		if matches = re.Match(`\t(\[.+\])\t`, entry); len(matches) > 1 {
			n.Flags = matches[1]
			v = true
		}

		if matches = re.Match(`\t([^\t]{1,32})[\n]`, entry); len(matches) > 1 {
			n.Ssid = matches[1]
			v = true
		}

		if v {
			answer = append(answer, &n)
		}
	}

	if len(answer) > 0 {
		return answer
	}

	return nil
}

func parseCommandBlock(block string) string {
	if matches := re.Match(`^([^\s]+)`, block); len(matches) > 1 {
		if matches[1] == "FAIL" {
			return "FAIL"
		}
		return matches[1]
	}
	return ""
}

func Status(iface string) (*status, error) {
	out, err := Wpa("-i", iface, "status")
	if err != nil {
		return nil, err
	}
	return parseStatusBlock(strings.Trim(out, " \r\n")), nil
}

func Bssid(iface string, ap string, ssid string) (string, error) {
	out, err := Wpa("-i", iface, "bssid", ssid, ap)
	if err != nil {
		return "", err
	}
	return parseCommandBlock(strings.Trim(out, " \r\n")), nil
}

func Reassociate(iface string) (string, error) {
	out, err := Wpa("-i", iface, "reassociate")
	if err != nil {
		return "", err
	}
	return parseCommandBlock(strings.Trim(out, " \r\n")), nil
}

func Set(iface string, name string, value interface{}) (string, error) {
	out, err := Wpa("-i", iface, "set", name, str.ToString(value))
	if err != nil {
		return "", err
	}
	return parseCommandBlock(strings.Trim(out, " \r\n")), nil
}

func AddNetwork(iface string) (string, error) {
	out, err := Wpa("-i", iface, "add_network")
	if err != nil {
		return "", err
	}
	return parseCommandBlock(strings.Trim(out, " \r\n")), nil
}

func SetNetwork(iface string, id int, name string, value string) (string, error) {
	out, err := Wpa("-i", iface, "set_network", str.ToString(id), name, value)
	if err != nil {
		return "", err
	}
	return parseCommandBlock(strings.Trim(out, " \r\n")), nil
}

func EnableNetwork(iface string, id int) (string, error) {
	out, err := Wpa("-i", iface, "enable_network", str.ToString(id))
	if err != nil {
		return "", err
	}
	return parseCommandBlock(strings.Trim(out, " \r\n")), nil
}

func DisableNetwork(iface string, id int) (string, error) {
	out, err := Wpa("-i", iface, "disable_network", str.ToString(id))
	if err != nil {
		return "", err
	}
	return parseCommandBlock(strings.Trim(out, " \r\n")), nil
}

func RemoveNetwork(iface string, id int) (string, error) {
	out, err := Wpa("-i", iface, "remove_network", str.ToString(id))
	if err != nil {
		return "", err
	}
	return parseCommandBlock(strings.Trim(out, " \r\n")), nil
}

func SelectNetwork(iface string, id int) (string, error) {
	out, err := Wpa("-i", iface, "select_network", str.ToString(id))
	if err != nil {
		return "", err
	}
	return parseCommandBlock(strings.Trim(out, " \r\n")), nil
}

func Scan(iface string) (string, error) {
	out, err := Wpa("-i", iface, "scan")
	if err != nil {
		return "", err
	}
	return parseCommandBlock(strings.Trim(out, " \r\n")), nil
}

func ScanResults(iface string) ([]*network, error) {
	out, err := Wpa("-i", iface, "scan_results")
	if err != nil {
		return nil, err
	}
	return parseScanResultsBlock(strings.Trim(out, " \r\n")), nil
}

func SaveConfig(iface string) (string, error) {
	out, err := Wpa("-i", iface, "save_config")
	if err != nil {
		return "", err
	}
	return parseCommandBlock(strings.Trim(out, " \r\n")), nil
}
