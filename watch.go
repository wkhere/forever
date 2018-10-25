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
	maxTick = 500 * time.Millisecond
)

// above will be configurable

func loop(w *fsnotify.Watcher) {

	type status uint
	const (
		stProcessed status = iota + 1
		stMinTick
		stMaxTick
	)

	var (
		ignoring, processing bool
		t0                   time.Time
		minTicker, maxTicker *time.Timer
		statusc              = make(chan status)
	)

	for {
		select {
		case ev, ok := <-w.Events:
			if !ok {
				return
			}
			e := evcatch{ev, time.Now()}
			debugf("* event i=%v p=%v", ignoring, processing)

			if ignoring {
				log("ignore event\t", e)
				continue
			}

			if processing {
				// for now, do nothing, TBD: queue the event
				log("processing did not end, skip event\t", e)
				continue
			}

			// if mintick comes during processing then nop
			// if mintick comes after processing then ignoring stops
			// processing end does not stop ignoring

			// max(processing, mintick) stops ignoring
			// min(processing, maxtick) also stops ignoring
			// if new request comes during processing but after maxtick then
			// processing will repeat - TBD

			ignoring = true
			processing = true
			t0 = time.Now()

			minTicker = time.AfterFunc(minTick, func() {
				debugf("* mintk i=%v p=%v", ignoring, processing)
				statusc <- stMinTick
			})
			maxTicker = time.AfterFunc(maxTick, func() {
				debugf("* maxtk i=%v p=%v", ignoring, processing)
				statusc <- stMaxTick
			})

			go func() {
				debugf("* proc! i=%v p=%v", ignoring, processing)
				process(e)
				statusc <- stProcessed
			}()

		case st := <-statusc:
			t1 := time.Now()

			debugf("* strcv i=%v p=%v st=%v", ignoring, processing, st)
			switch st {
			case stProcessed:
				processing = false

				if t1.Sub(t0) < minTick {
					// nothing to do more, minTick will come
					continue
				}
				minTicker.Stop() // in case of some subtle scheduling slip
				maxTicker.Stop()
				ignoring = false

			case stMinTick:
				if !processing {
					ignoring = false
					maxTicker.Stop()
				}

			case stMaxTick:
				if processing && ignoring {
					ignoring = false
				}
			}

		case err, ok := <-w.Errors:
			if !ok {
				return
			}
			log("received error:", err)
		}
	}
}
