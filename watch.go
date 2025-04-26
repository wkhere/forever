package main

import (
	"fmt"
	"os"
	"time"

	"github.com/fsnotify/fsnotify"
)

type watcher struct {
	*fsnotify.Watcher
	dirs     []string
	timeslot time.Duration
}

func newWatcher(t time.Duration) (w *watcher, err error) {
	w = &watcher{timeslot: t}
	w.Watcher, err = fsnotify.NewWatcher()
	return
}

type statusCode uint

const (
	stRunFinished statusCode = iota + 1
	stTimeslotEnded
)

type status struct {
	statusCode
	t0 time.Time
}

type runCause int

const (
	runFirst runCause = iota + 1
	runAwakened
)

func loop(w *watcher, p *prog) {

	var (
		ignoring, running bool
		timer             *time.Timer
		statusc           = make(chan status)
	)

	startRunning := func(why runCause) {
		ignoring = true
		running = true

		t0 := time.Now()

		var causeText string
		switch why {
		case runAwakened:
			causeText = "awakened"
		case runFirst:
			causeText = "started"
		}

		logfBlue("[forever %s %s]", causeText, timef(t0))

		timer = time.AfterFunc(w.timeslot, func() {
			watchdebug("tslot i=%v r=%v", ignoring, running)
			statusc <- status{stTimeslotEnded, t0}
		})

		go func() {
			watchdebug("run-> i=%v r=%v", ignoring, running)

			pst, err := p.run()
			switch {
			case pst == nil:
				logfRed("%v", err)
			default:
				t := time.Now()
				if err != nil {
					logfRed("run `%v` failed: %v", p, err)
				}
				logfBlue("[%s]", pstatef(pst, t.Sub(t0)))
			}
			statusc <- status{stRunFinished, t0}
		}()
	}

	watchdebug("will start for the first time")
	startRunning(runFirst)

	for {
		select {
		case ev, ok := <-w.Events:
			if !ok {
				return
			}
			watchdebug("got event %v, i=%v r=%v", ev, ignoring, running)

			// The timeslot covers a series of create/write events related to
			// the same file-change operation.
			//
			// * if the timeslot ends during running then nop
			// * if the timeslot ends after a finished run then ignoring stops
			// * finished run doesn't stop ignoring unless timeslot ended

			// so, max(running_time, timeslot) stops ignoring

			// if a new request comes during running, is ignored as a conse-
			// quence of the scenario above

			if ignoring {
				watchdebug("ignore event %v", ev)
				continue
			}

			if running {
				panic("watch: running while not ignoring shouldn't happen")
			}

			watchdebug("will process event %v", ev)
			startRunning(runAwakened)

		case st := <-statusc:
			t1 := time.Now()

			watchdebug("strcv i=%v r=%v st=%v", ignoring, running, st)

			switch st.statusCode {
			case stRunFinished:
				running = false

				if t1.Sub(st.t0) < w.timeslot {
					// nothing to do more, timeslot will end
					continue
				}
				timer.Stop() // in case of some subtle scheduling slip
				ignoring = false

			case stTimeslotEnded:
				if !running {
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
	case stRunFinished:
		s = "stRunFinished"
	case stTimeslotEnded:
		s = "stTimeslotEnded"
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
