package main // import "github.com/wkhere/forever"

import (
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/fsnotify/fsnotify"
)

func init() {
	flag.Parse()
}

func main() {
	var dir string

	switch a := flag.Args(); len(a) {
	case 0:
		dir = "."
	case 1:
		dir = a[0]
	default:
		log("give only 1 directory as an argument")
		os.Exit(2)
	}

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

	// watcher should add all files before chdir
	feedWatcher(w, dir)

	err = os.Chdir(dir)
	if err != nil {
		log("cannot prepare:", err)
		os.Exit(1)
	}

	go loop(w)

	done := make(chan struct{})
	<-done
}

func log(msgs ...interface{}) {
	fmt.Fprintln(os.Stderr, msgs...)
}

func logf(format string, msgs ...interface{}) {
	log(fmt.Sprintf(format, msgs...))
}

func logBlue(s string) {
	io.WriteString(os.Stderr, "\033[34m")
	io.WriteString(os.Stderr, s)
	io.WriteString(os.Stderr, "\033[0m\n")
}
