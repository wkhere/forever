package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/wkhere/redbuffer"
)

type progConfigT struct {
	explicitProg bool
	prog         string
	args         []string
	redOnError   bool
}

const stepfile = ".forever.step"

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
	return run(pc.prog, pc.args, pc.redOnError)
}

var defaultProgsDescription = fmt.Sprintf(
	`
	sh %s, if that file exists
	make
`,
	stepfile,
)

func (pc *progConfigT) processDefaultProgs() (*os.ProcessState, error) {

	switch _, err := os.Stat(stepfile); {
	case err == nil:
		return run("sh", []string{"-e", stepfile}, pc.redOnError)

	case errors.Is(err, os.ErrNotExist):
		break

	case err != nil:
		return nil, fmt.Errorf("Unexpected error when looking for %s: %s",
			stepfile, err)
	}

	return run("make", nil, pc.redOnError)
}

func run(p string, args []string, redOnError bool) (*os.ProcessState, error) {
	c := exec.Command(p, args...)
	c.Stdout = os.Stdout
	w := redbuffer.New(os.Stderr)
	c.Stderr = w
	err := c.Run()
	w.FlushInRed(redOnError && err != nil)
	return c.ProcessState, err
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
