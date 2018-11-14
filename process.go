package main

import (
	"errors"
	"os"
	"os/exec"
)

var progs = []string{
	"./.forever.step",
	"make",
}

func process() (*os.ProcessState, error) {
	for _, p := range progs {
		if _, err := exec.LookPath(p); err != nil {
			continue
		}
		c := exec.Command(p)
		c.Stdout = os.Stdout
		c.Stderr = os.Stderr
		err := c.Run()
		return c.ProcessState, err
	}
	return nil, errProgNotRun
}

var errProgNotRun = errors.New("could not run suitable program")
