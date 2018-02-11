package osencap

import (
	"os/exec"
	"log"
	"io/ioutil"
)

var Exec = func(command string, args ...string) (string, error) {
	cmdpath, err := exec.LookPath(command)
	if err != nil {
		log.Fatalf("exec.LookPath err: %v, cmd: %s", err, command)
	}

	cmd := exec.Command(cmdpath, args...)
	stdout, err := cmd.StdoutPipe()
	if err := cmd.Start(); err != nil {
		log.Fatal(err)
	}

	out, _ := ioutil.ReadAll(stdout)

	if err := cmd.Wait(); err != nil {
		log.Fatal(err)
	}

	return string(out), nil
}
