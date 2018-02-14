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

	c := exec.Command(cmdpath, args...)
	stdout, err := c.StdoutPipe()
	if err != nil {
		return "", err
	}
	defer stdout.Close()

	if err := c.Start(); err != nil {
		return "", err
	}

	out, err := ioutil.ReadAll(stdout)
	if err != nil {
		return "", err
	}

	if err := c.Wait(); err != nil {
		return "", err
	}

	return string(out), nil
}
