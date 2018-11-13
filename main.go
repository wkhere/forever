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
	var dirs []string

	switch a := flag.Args(); {
	case len(a) == 0:
		dirs = []string{"."}
	default:
		dirs = a
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

	go loop(w)

	for _, dir := range dirs {
		feedWatcher(w, dir)
	}

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
