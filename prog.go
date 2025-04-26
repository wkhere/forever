package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

type prog struct {
	path     string
	args     []string
	explicit bool
}

const stepfile = ".forever.step"

func (p *prog) run() (*os.ProcessState, error) {
	if !p.explicit {
		return p.runDefaultProgs()
	}
	return p.runProg()
}

func (p *prog) runProg() (*os.ProcessState, error) {
	if _, err := exec.LookPath(p.path); err != nil {
		return nil, fmt.Errorf("could not run given program: %v", err)
	}
	return run(p.path, p.args)
}

var defaultProgsDescription = fmt.Sprintf(
	`
	sh %s, if that file exists
	make
`,
	stepfile,
)

func (p *prog) runDefaultProgs() (*os.ProcessState, error) {

	switch _, err := os.Stat(stepfile); {
	case err == nil:
		return run("sh", []string{"-e", stepfile})

	case errors.Is(err, os.ErrNotExist):
		break

	case err != nil:
		return nil, fmt.Errorf("Unexpected error when looking for %s: %s",
			stepfile, err)
	}

	return run("make", nil)
}

func run(p string, args []string) (*os.ProcessState, error) {
	c := exec.Command(p, args...)
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	err := c.Run()
	return c.ProcessState, err
}

// is this dead code? :
func (p prog) String() (s string) {
	if !p.explicit {
		s = "(default) "
	}
	switch {
	case len(p.args) == 0:
		s += p.path
	case len(p.args) > 4:
		return p.path + " ..."
	default:
		s += p.path + " " + strings.Join(p.args, " ")
	}
	return
}
