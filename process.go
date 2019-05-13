package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

type progConfigT struct {
	explicitProg bool
	prog         string
	args         []string
}

func (pc progConfigT) String() (s string) {
	if !pc.explicitProg {
		s = "(default) "
	}
	switch {
	case len(pc.args) == 0:
		s += pc.prog
	case len(pc.args) > 4:
		return pc.prog + " ..."
	default:
		s += pc.prog + " " + strings.Join(pc.args, " ")
	}
	return
}

var defaultProgs = []string{
	"./.forever.step",
	"make",
}

func (pc *progConfigT) process() (*os.ProcessState, error) {
	if !pc.explicitProg {
		chosen, ps, err := processDefaultProgs()
		pc.prog = chosen
		return ps, err
	}
	return processProg(pc.prog, pc.args)
}

func processProg(p string, args []string) (*os.ProcessState, error) {
	if _, err := exec.LookPath(p); err != nil {
		return nil, fmt.Errorf("could not run given program: %v", err)
	}
	return run(p, args)
}

func processDefaultProgs() (string, *os.ProcessState, error) {
	for _, p := range defaultProgs {
		if _, err := exec.LookPath(p); err != nil {
			continue
		}
		ps, err := run(p, nil)
		return p, ps, err
	}
	return "", nil, fmt.Errorf("could not run any of default programs")
}

func run(p string, args []string) (*os.ProcessState, error) {
	c := exec.Command(p, args...)
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	err := c.Run()
	return c.ProcessState, err
}
