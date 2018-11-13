package main

import (
	"os"
	"os/exec"
)

var progs = []string{
	"./.forever.step",
	"make",
}

func process() {
	var done bool
	for _, p := range progs {
		if done {
			break
		}
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
		done = true
	}
	if !done {
		log("could not find suitable program to run")
	}
}

func progNotFound(err error) bool {
	return os.IsNotExist(err) || err == exec.ErrNotFound
}
