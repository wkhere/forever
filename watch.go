package main // import "github.com/wkhere/forever"

import (
	"fmt"
	"time"

	"github.com/fsnotify/fsnotify"
)

type evcatch struct {
	ev fsnotify.Event
	t  time.Time
}

func (e evcatch) String() string {
	return fmt.Sprintf("{%v %v}", e.t, e.ev)
}

var (
	minTick = 200 * time.Millisecond
	// tick time will be configurable
)

func loop(w *fsnotify.Watcher) {

	type status uint
	const (
		stProcessed status = iota + 1
		stMinTick
	)

	var (
		ignoring, processing bool
		t0                   time.Time
		minTicker            *time.Timer
		statusc              = make(chan status)
	)

	startProcessing := func() {
		ignoring = true
		processing = true
		t0 = time.Now()

		minTicker = time.AfterFunc(minTick, func() {
			debugf("watch: mintk i=%v p=%v", ignoring, processing)
			statusc <- stMinTick
		})

		go func() {
			debugf("watch: proc! i=%v p=%v", ignoring, processing)
			process()
			statusc <- stProcessed
		}()
	}

	startProcessing()
	logBlue("[forever started]")

	for {
		select {
		case ev, ok := <-w.Events:
			if !ok {
				return
			}
			e := evcatch{ev, time.Now()}
			debugf("watch: event i=%v p=%v", ignoring, processing)

			if ignoring {
				debugf("watch: ignore\t%s", e)
				continue
			}

			if processing {
				panic("watch: processing while not ignoring shouldn't happen")
			}

			// if mintick comes during processing then nop
			// if mintick comes after processing then ignoring stops
			// processing end doesn't stop ignoring unless mintick time passed

			// so, max(processing, mintick) stops ignoring

			// if new request comes during processing, is ignored as a conse-
			// quence of the scenario above

			debugf("watch: process\t%s", e)
			startProcessing()
			logBlue("[forever awakened]")

		case st := <-statusc:
			t1 := time.Now()

			debugf("watch: strcv i=%v p=%v st=%v", ignoring, processing, st)
			switch st {
			case stProcessed:
				processing = false

				if t1.Sub(t0) < minTick {
					// nothing to do more, minTick will come
					continue
				}
				minTicker.Stop() // in case of some subtle scheduling slip
				ignoring = false

			case stMinTick:
				if !processing {
					ignoring = false
				}
			}

		case err, ok := <-w.Errors:
			if !ok {
				return
			}
			debugf("watch: received error:", err)
		}
	}
}
