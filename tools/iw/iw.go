package iw

import (
	"github.com/taoyuan/go-awt/infra/osencap"
	"github.com/taoyuan/go-awt/tools/iu"
	"regexp"
	"sort"
	"github.com/thoas/go-funk"
	"strings"
	"strconv"
	"github.com/taoyuan/go-awt/infra/re"
)

type Cell struct {
	Ssid       string
	Address    string
	Frequency  int
	Signal     float64
	LastSeenMs int
	Channel    int
	Security   string
}

type Cells []*Cell

func (s Cells) Len() int      { return len(s) }
func (s Cells) Swap(i, j int) { s[i], s[j] = s[j], s[i] }

type BySignal struct{ Cells }

func (s BySignal) Less(i, j int) bool { return s.Cells[i].Signal < s.Cells[j].Signal }

/**
 *
 */
var Iw = func(args ...string) (string, error) {
	return osencap.Exec("iw", args...)
}

func parseCell(data string) (*Cell) {
	cell := Cell{}

	var matches []string

	if matches = re.Match(`BSS ([0-9A-Fa-f:-]{17})\(on`, data); len(matches) > 0 {
		cell.Address = matches[1]
	}

	if matches = re.Match(`freq: ([0-9]+)`, data); len(matches) > 0 {
		cell.Frequency, _ = strconv.Atoi(matches[1])
	}

	if matches = re.Match(`signal: (-?[0-9.]+) dBm`, data); len(matches) > 0 {
		cell.Signal, _ = strconv.ParseFloat(matches[1], 64)
	}

	if matches = re.Match(`last seen: ([0-9]+) ms ago`, data); len(matches) > 0 {
		cell.LastSeenMs, _ = strconv.Atoi(matches[1])
	}

	if matches = re.Match(`SSID: \\x00`, data); len(matches) > 0 {
		//
	} else if matches = re.Match(`SSID: ([^\n]*)`, data); len(matches) > 0 {
		cell.Ssid = matches[1]
	}

	if matches = re.Match(`DS Parameter set: channel ([0-9]+)`, data); len(matches) > 0 {
		cell.Channel, _ = strconv.Atoi(matches[1])
	} else if matches = re.Match(`\* primary channel: ([0-9]+)`, data); len(matches) > 0 {
		cell.Channel, _ = strconv.Atoi(matches[1])
	}

	if matches = re.Match(`RSN:[\s*]+Version: 1`, data); len(matches) > 0 {
		cell.Security = "wpa2"
	} else if matches = re.Match(`WPA:[\s*]+Version: 1`, data); len(matches) > 0 {
		cell.Security = "wpa"
	} else if matches = re.Match(`capability: ESS Privacy`, data); len(matches) > 0 {
		cell.Security = "wep"
	} else if matches = re.Match(`capability: ESS`, data); len(matches) > 0 {
		cell.Security = "open"
	}

	return &cell
}

func parseScan(showHidden bool, data string) Cells {
	reg := regexp.MustCompile(`(^|\n)(BSS )`)

	parts := reg.Split(data, -1)
	parts = funk.Map(parts, func(s string) string { return strings.Trim(s, " ") }).([]string)
	parts = funk.Filter(parts, func(s string) bool { return s != "" }).([]string)
	parts = funk.Map(parts, func(s string) string { return "BSS " + s }).([]string)

	cells := Cells{}
	for _, part := range parts {
		cell := parseCell(part)
		cells = append(cells, cell)
	}

	if showHidden {
		cells = filterCells(cells, hasAddr)
	} else {
		cells = filterCells(cells, hasSsid)
	}

	sort.Sort(BySignal{cells})

	return cells
}

func filterCells(cells Cells, f func(*Cell) bool) Cells {
	vsf := make(Cells, 0)
	for _, cell := range cells {
		if f(cell) {
			vsf = append(vsf, cell)
		}
	}
	return vsf
}

func hasAddr(cell *Cell) bool {
	return cell.Address != ""
}

func hasSsid(cell *Cell) bool {
	return cell.Ssid != ""
}

func Scan(iface string) (Cells, error) {
	iface, _ = iu.ResolveIface(iface)
	out, err := Iw("dev", iface, "scan")
	if err != nil {
		return nil, err
	}
	cells := parseScan(true, out)
	return cells, nil
}
