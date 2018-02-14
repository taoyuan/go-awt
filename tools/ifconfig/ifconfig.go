package ifconfig

import (
	"go-awt/infra/osencap"
	"strings"
	"go-awt/infra/re"
	"regexp"
)

type status struct {
	Iface          string
	Link           string
	Address        string
	Ipv6Address    string
	Ipv4Address    string
	Ipv4Broadcast  string
	Ipv4SubnetMask string
	Up             bool
	Broadcast      bool
	Running        bool
	Multicast      bool
	Loopback       bool
}

func parseStatusBlock(block string) *status {
	st := status{}
	var matches []string

	if matches = re.Match(`^([^\s]+)`, block); len(matches) > 1 {
		st.Iface = strings.ToLower(matches[1])
	}

	if matches = re.Match(`Link encap:\s*([^\s]+)`, block); len(matches) > 1 {
		st.Link = strings.ToLower(matches[1])
	}

	if matches = re.Match(`HWaddr\s+([^\s]+)`, block); len(matches) > 1 {
		st.Address = strings.ToLower(matches[1])
	}

	if matches = re.Match(`inet6\s+addr:\s*([^\s]+)`, block); len(matches) > 1 {
		st.Ipv6Address = matches[1]
	}

	if matches = re.Match(`inet\s+addr:\s*([^\s]+)`, block); len(matches) > 1 {
		st.Ipv4Address = matches[1]
	}

	if matches = re.Match(`Bcast:\s*([^\s]+)`, block); len(matches) > 1 {
		st.Ipv4Broadcast = matches[1]
	}

	if matches = re.Match(`Mask:\s*([^\s]+)`, block); len(matches) > 1 {
		st.Ipv4SubnetMask = matches[1]
	}

	if ret, _ := regexp.Match(`UP`, []byte(block)); ret {
		st.Up = true
	}

	if ret, _ := regexp.Match(`BROADCAST`, []byte(block)); ret {
		st.Broadcast = true
	}

	if ret, _ := regexp.Match(`RUNNING`, []byte(block)); ret {
		st.Running = true
	}

	if ret, _ := regexp.Match(`MULTICAST`, []byte(block)); ret {
		st.Multicast = true
	}

	if ret, _ := regexp.Match(`LOOPBACK`, []byte(block)); ret {
		st.Loopback = true
	}

	return &st
}

func parseStatus(stdout string) []*status {
	blocks := strings.Split(stdout, "\n\n")
	var answer []*status
	for _, block := range blocks {
		answer = append(answer, parseStatusBlock(block))
	}
	return answer
}

var Ifconfig = func(args ...string) (string, error) {
	return osencap.Exec("ifconfig", args...)
}

func Status(iface string) ([]*status, error) {
	var out string
	var err error
	if iface != "" {
		out, err = Ifconfig(iface)
	} else {
		out, err = Ifconfig("-a")
	}

	if err != nil {
		return nil, err
	}
	return parseStatus(strings.Trim(out, " \r\n")), nil
}
