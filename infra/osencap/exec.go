package osencap

import (
	"os/exec"
	"bytes"
	"fmt"
	"go-awt/infra"
)

var Exec = func(command string, args ...string) (string, error) {
	cmdpath, err := exec.LookPath(command)
	if err != nil {
		fmt.Errorf("exec.LookPath err: %v, cmd: %s", err, command)
		return "", infra.ErrExecLookPathFailed
	}

	cmd := exec.Command(cmdpath, args...)
	output, err := cmd.StdoutPipe()
	if err := cmd.Start(); err != nil {
		fmt.Errorf("cmd.Start err: %v, cmd: %s", err, command)
	}

	buf := new(bytes.Buffer)
	buf.ReadFrom(output)
	return buf.String(), nil
}
