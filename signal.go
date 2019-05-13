// +build !windows

package main

import (
	"os"
	"os/signal"
	"syscall"
)

func (w *watcher) installSignal() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP)

	go func() {
		for {
			<-c
			logf("watching %d dir(s):", len(w.dirs))
			for _, d := range w.dirs {
				log("\t", d)
			}
		}
	}()
}
