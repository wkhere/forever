package main

import (
	"time"

	"github.com/fsnotify/fsnotify"
)

type watcher struct {
	*fsnotify.Watcher
	dirs    []string
	minTick time.Duration
}

func newWatcher(t time.Duration) (w *watcher, err error) {
	w = &watcher{minTick: t}
	w.Watcher, err = fsnotify.NewWatcher()
	return
}
