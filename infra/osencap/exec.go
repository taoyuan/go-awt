package osencap

import (
	"fmt"
	"os/exec"
	"go-awt/infra"
	"log"
	"bytes"
)

var Exec = func(cmd string, args ...string) (string, error) {
	cmdpath, err := exec.LookPath(cmd)
	if err != nil {
		fmt.Errorf("exec.LookPath err: %v, cmd: %s", err, cmd)
		return "", infra.ErrExecLookPathFailed
	}

	command := exec.Command(cmdpath, args...)
	output, err := command.StdoutPipe()
	if err := command.Start(); err != nil {
		log.Fatal(err)
	}

	buf := new(bytes.Buffer)
	buf.ReadFrom(output)
	return buf.String(), nil
	//
	//if err != nil {
	//	fmt.Errorf("exec.Command.CombinedOutput err: %v, cmd: %s", err, cmd)
	//	return "", infra.ErrExecCombinedOutputFailed
	//}
	////fmt.Println("CMD[", cmdpath, "]ARGS[", args, "]OUT[", string(output), "]")
	//return string(output), nil
}
