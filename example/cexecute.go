package main

import (
	"github.com/go-cmd/cmd"
	"fmt"
	"time"
)

func main() {
	// Start a long-running process, capture stdout and stderr
	c := cmd.NewCmd("node", "/users/taoyuan/temp/test.js")
	statusChan := c.Start()

	// Print last line of stdout every 2s
	go func() {
		for range time.NewTicker(2 * time.Second).C {
			status := c.Status()
			//n := len(status.Stdout)
			fmt.Println(status.Stdout)
		}
	}()

	// Stop command after 1 hour
	go func() {
		<-time.After(1 * time.Minute)
		c.Stop()
	}()

	// Check if command is done
	select {
	case <-statusChan:
	}

	// Block waiting for command to exit, be stopped, or be killed
	<-statusChan
}
