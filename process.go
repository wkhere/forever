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
	redbuf       bool
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
		return pc.processDefaultProgs()
	}
	return pc.processProg()
}

func (pc *progConfigT) processProg() (*os.ProcessState, error) {
	if _, err := exec.LookPath(pc.prog); err != nil {
		return nil, fmt.Errorf("could not run given program: %v", err)
	}
	return run(pc.prog, pc.args, pc.redbuf)
}

func (pc *progConfigT) processDefaultProgs() (*os.ProcessState, error) {
	for _, p := range defaultProgs {
		if _, err := exec.LookPath(p); err != nil {
			continue
		}
		pc.prog = p
		ps, err := run(p, nil, pc.redbuf)
		return ps, err
	}
	return nil, fmt.Errorf("could not run any of default programs")
}

func run(p string, args []string, redbuf bool) (*os.ProcessState, error) {
	c := exec.Command(p, args...)
	c.Stdout = os.Stdout
	w := newRedbufWriter(os.Stderr)
	c.Stderr = w
	err := c.Run()
	w.FlushInRed(redbuf && err != nil)
	return c.ProcessState, err
}
