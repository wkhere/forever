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
		c := exec.Command(p)
		c.Stdout = os.Stdout
		c.Stderr = os.Stderr
		err := c.Run()
		if progNotFound(err) {
			continue
		}
		if err != nil {
			log(err)
			continue
		}

		return c.ProcessState, nil
	}

	return nil, errProgNotRun
}

var errProgNotRun = errors.New("could not run suitable program")

func progNotFound(err error) bool {
	return os.IsNotExist(err) || err == exec.ErrNotFound
}
