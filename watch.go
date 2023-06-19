package main

import (
	"fmt"
	"os"
	"time"
)

type statusCode uint

const (
	stProcessed statusCode = iota + 1
	stMinTick
)

type status struct {
	statusCode
	t0 time.Time
}

type processingCause int

const (
	procStarted processingCause = iota + 1
	procAwakened
)

func loop(w *watcher, pc *progConfigT) {

	var (
		ignoring, processing bool
		minTicker            *time.Timer
		statusc              = make(chan status)
	)

	startProcessing := func(why processingCause) {
		ignoring = true
		processing = true

		t0 := time.Now()

		var causeText string
		switch why {
		case procAwakened:
			causeText = "awakened"
		case procStarted:
			causeText = "started"
		}

		logfBlue("[forever %s %s]", causeText, timef(t0))

		minTicker = time.AfterFunc(w.minTick, func() {
			watchdebug("mintk i=%v p=%v", ignoring, processing)
			statusc <- status{stMinTick, t0}
		})

		go func() {
			watchdebug("proc! i=%v p=%v", ignoring, processing)
			pst, err := pc.process()
			switch {
			case pst == nil:
				logfRed("%v", err)
			default:
				t := time.Now()
				if err != nil {
					logfRed("process `%v` failed: %v", pc, err)
				}
				logfBlue("[%s]", pstatef(pst, t.Sub(t0)))
			}
			statusc <- status{stProcessed, t0}
		}()
	}

	watchdebug("will start for first time")
	startProcessing(procStarted)

	for {
		select {
		case ev, ok := <-w.Events:
			if !ok {
				return
			}
			watchdebug("got event %v, i=%v p=%v", ev, ignoring, processing)

			if ignoring {
				watchdebug("ignore event %v", ev)
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

			watchdebug("will process event %v", ev)
			startProcessing(procAwakened)

		case st := <-statusc:
			t1 := time.Now()

			watchdebug("strcv i=%v p=%v st=%v", ignoring, processing, st)

			switch st.statusCode {
			case stProcessed:
				processing = false

				if t1.Sub(st.t0) < w.minTick {
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
			log("watch: received error:", err)
		}
	}
}

func (c statusCode) String() (s string) {
	switch c {
	case stProcessed:
		s = "stProcessed"
	case stMinTick:
		s = "stMinTick"
	}
	return
}
func (s status) String() string {
	return fmt.Sprintf("{code=%s t0=%s}", s.statusCode, timef_ns(s.t0))
}

func timef(t time.Time) string {
	return t.Format("15:04:05")
}

func timef_ns(t time.Time) string {
	return t.Format("15:04:05.000000000")
}

func watchdebug(format string, a ...interface{}) {
	debugf(fmt.Sprintf("watch at %s: ", timef_ns(time.Now()))+format,
		a...)
}

func pstatef(pst *os.ProcessState, wall time.Duration) string {
	const ms = time.Millisecond
	var (
		sys  = pst.SystemTime()
		user = pst.UserTime()
		pcpu = float64(user+sys) / float64(wall)
	)
	if mem, ok := rusageExtras.getMemStats(pst); ok {
		maxrssMB := float64(mem.maxRss) / 1024
		return fmt.Sprintf(
			"%s user  %s sys  %.2f%% cpu  %s total,  %.1fM rss  %d/%d flt",
			user.Round(ms), sys.Round(ms), pcpu*100, wall.Round(ms), maxrssMB,
			mem.minFlt, mem.majFlt,
		)
	}
	return fmt.Sprintf("%s user  %s sys  %.2f%% cpu  %s total",
		user.Round(ms), sys.Round(ms), pcpu*100, wall.Round(ms))
}
