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

func (pc progConfigT) String() string {
	if !pc.explicitProg {
		return "<undefined program, will use defaults>"
	}
	if len(pc.args) == 0 {
		return pc.prog
	}
	if len(pc.args) > 4 {
		return pc.prog + " ..."
	}
	return pc.prog + " " + strings.Join(pc.args, " ")

}

var defaultProgs = []string{
	"./.forever.step",
	"make",
}

func (pc *progConfigT) process() (*os.ProcessState, error) {
	if !pc.explicitProg {
		return processDefaultProgs()
	}
	return processProg(pc.prog, pc.args)
}

func processProg(p string, args []string) (*os.ProcessState, error) {
	if _, err := exec.LookPath(p); err != nil {
		return nil, fmt.Errorf("could not run given program: %v", err)
	}
	return run(p, args)
}

func processDefaultProgs() (*os.ProcessState, error) {
	for _, p := range defaultProgs {
		if _, err := exec.LookPath(p); err != nil {
			continue
		}
		return run(p, nil)
	}
	return nil, fmt.Errorf("could not run any of default program")
}

func run(p string, args []string) (*os.ProcessState, error) {
	c := exec.Command(p, args...)
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	err := c.Run()
	return c.ProcessState, err
}
