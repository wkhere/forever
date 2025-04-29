package main // import "github.com/wkhere/forever"

import (
	"fmt"
	"io"
	"os"
	"time"
)

type config struct {
	dir      string
	timeslot time.Duration

	prog prog

	verbose bool
	help    func(io.Writer)
}

func defaults() *config {
	return &config{dir: ".", timeslot: 200 * time.Millisecond}
}

func main() {
	config, err := parseArgs(os.Args[1:])
	if err != nil {
		die(2, err)
	}
	if config.help != nil {
		config.help(os.Stdout)
		os.Exit(0)
	}

	err = run(config)
	if err != nil {
		die(1, err)
	}
}

func die(exitcode int, err error) {
	fmt.Fprintln(os.Stderr, err)
	os.Exit(exitcode)
}

func run(c *config) error {

	err := os.Chdir(c.dir)
	if err != nil {
		return err
	}

	w, err := newWatcher(c.timeslot)
	if err != nil {
		return err
	}

	// watcher should add all files before looping
	err = w.feed()
	if err != nil {
		return err
	}
	w.installSignal()

	go loop(w, &c.prog)
	select {} // block forever
}
