package main // import "github.com/wkhere/forever"

import (
	"flag"
	"fmt"
	"os"

	"github.com/fsnotify/fsnotify"
)

func init() {
	flag.Parse()
}

func main() {

	w, err := fsnotify.NewWatcher()
	if err != nil {
		log("cannot start watcher: ", err)
		os.Exit(1)
	}
	defer func() {
		err = w.Close()
		if err != nil {
			log("error during closing watcher:", err)
		}
	}()

	go loop(w)

	for _, f := range flag.Args() {
		err := w.Add(f)
		if err != nil {
			log("skipping", f, "error:", err)
		}
	}

	done := make(chan struct{})
	<-done
}

func log(msgs ...interface{}) {
	fmt.Fprintln(os.Stderr, msgs...)
}
