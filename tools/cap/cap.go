package cap

import (
	"github.com/go-cmd/cmd"
	"errors"
	"strconv"
)

var SSID = "AWT"
var CommandCreateAP = "create_ap"

type Options struct {
	Iface          string
	IfaceSharing   string
	Ssid           string
	Pass           string
	Bridge         bool
	Driver         string
	HtCapab        string
	Ieee80211n     bool
	IsolateClients bool
	NoVirt         bool
}

type AP struct {
	cmd    *cmd.Cmd
	ch     <-chan cmd.Status
	curr   int
}

func CreateAP(opts *Options) (*AP, error) {
	var args []string

	if opts.IfaceSharing == "" {
		args = append(args, "-n")
	}

	if opts.Bridge {
		args = append(args, "-m", "bridge")
	}

	if opts.Driver != "" {
		args = append(args, "--driver", opts.Driver)
	}

	if opts.Ieee80211n {
		args = append(args, "--ieee80211n")
	}

	if opts.HtCapab != "" {
		args = append(args, "--ht_capab", "'"+opts.HtCapab+"'")
	}

	if opts.IsolateClients {
		args = append(args, "--isolate-clients")
	}

	iface := opts.Iface
	if iface == "" {
		return nil, errors.New("iface can not be blank")
	}

	args = append(args, opts.Iface)

	if opts.IfaceSharing != "" {
		args = append(args, opts.IfaceSharing)
	}

	ssid := opts.Ssid
	if ssid == "" {
		ssid = SSID
	}

	args = append(args, ssid)

	if opts.Pass != "" {
		args = append(args, opts.Pass)
	}

	return &AP{
		cmd:    cmd.NewCmd(CommandCreateAP, args...),
		ch:     nil,
		curr:	0,
	}, nil
}

func (ap *AP) Start() <-chan cmd.Status {
	if ap.ch != nil {
		return ap.ch
	}
	ap.ch = ap.cmd.Start()
	ap.curr = 0
	return ap.ch
}

func (ap *AP) Stop() error {
	if ap.ch == nil {
		return errors.New("ap: not started")
	}
	return ap.cmd.Stop()
}

func (ap *AP) Wait() error {
	if ap.ch == nil {
		return errors.New("ap: not started")
	}
	status := <-ap.ch
	if status.Error != nil {
		return status.Error
	}
	if status.Exit != 0 {
		return errors.New("exit status " + strconv.Itoa(status.Exit))
	}
	return nil
}

func (ap *AP) Status() cmd.Status {
	return ap.cmd.Status()
}

func (ap *AP) Output() []string {
	status := ap.cmd.Status()
	if ap.curr < len(status.Stdout) {
		output := status.Stdout[ap.curr:]
		ap.curr = len(status.Stdout)
		return output
	}
	return nil
}
