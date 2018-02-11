package osencap

import (
	"os/exec"
	"io/ioutil"
)

var Exec = func(command string, args ...string) (string, error) {
	cmdpath, err := exec.LookPath(command)
	if err != nil {
		return "", err
	}

	cmd := exec.Command(cmdpath, args...)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return "", err
	}
	defer stdout.Close()

	if err := cmd.Start(); err != nil {
		return "", err
	}

	out, err := ioutil.ReadAll(stdout)
	if err != nil {
		return "", err
	}

	if err := cmd.Wait(); err != nil {
		return "", err
	}

	return string(out), nil
}
