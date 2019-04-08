package main

import (
	"fmt"
	"os"
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

func loop(w *fsnotify.Watcher, minTick time.Duration, pc *progConfigT) {

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

		minTicker = time.AfterFunc(minTick, func() {
			debugf("watch: mintk i=%v p=%v", ignoring, processing)
			statusc <- stMinTick
		})

		go func() {
			debugf("watch: proc! i=%v p=%v", ignoring, processing)
			pst, err := pc.process()
			switch {
			case pst == nil:
				log(err)
			default:
				t := time.Now()
				if err != nil {
					logf("process `%v` failed: %v", pc, err)
				}
				logBlue(fmt.Sprintf("[%s]", pstatef(pst, t.Sub(t0))))
			}
			statusc <- stProcessed
		}()
	}

	t0 = time.Now()
	logBlue(fmt.Sprintf("[forever started %s]", timef(t0)))
	startProcessing()

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
			t0 = time.Now()
			logBlue(fmt.Sprintf("[forever awakened %s]", timef(t0)))
			startProcessing()

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

func timef(t time.Time) string {
	return t.Format("15:04:05")
}

func pstatef(pst *os.ProcessState, wall time.Duration) string {
	const ms = time.Millisecond
	var (
		sys  = pst.SystemTime()
		user = pst.UserTime()
		pcpu = float64(user+sys) / float64(wall)
	)
	if maxrss, ok := rusageExtras.maxRss(pst); ok {
		return fmt.Sprintf("%s user  %s sys  %.2f%% cpu  %s total, rss %dk",
			user.Round(ms), sys.Round(ms), pcpu*100, wall.Round(ms), maxrss)
	}
	return fmt.Sprintf("%s user  %s sys  %.2f%% cpu  %s total",
		user.Round(ms), sys.Round(ms), pcpu*100, wall.Round(ms))
}
