// +build !windows

package main

import (
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

const writeDirsOnSignal = true

var writeDirsOutputPattern = "/tmp/forever-$PID"

func (w *watcher) installSignal() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP)

	go func() {
		for {
			<-c
			fn, err := writeDirs(w.dirs)
			if err != nil {
				log("failed to write list of watched dirs:", err)
				continue
			}
			log("list of watched dirs written to:", fn)
		}
	}()
}

func writeDirs(dirs []string) (fn string, err error) {
	fn = strings.Replace(
		writeDirsOutputPattern, "$PID", fmt.Sprintf("%d", os.Getpid()), 1,
	)
	f, err := os.Create(fn)
	if err != nil {
		return "", err
	}
	_, err = fmt.Fprintf(f, "watching %d dir(s):\n", len(dirs))
	if err != nil {
		return "", err
	}
	for _, d := range dirs {
		_, err = fmt.Fprintf(f, "\t%s\n", d)
		if err != nil {
			return "", err
		}
	}
	if err = f.Close(); err != nil {
		return "", err
	}
	return fn, nil
}
