package main

import (
	"fmt"
	"os"
	"time"

	"github.com/fsnotify/fsnotify"
)

type watcher struct {
	*fsnotify.Watcher
	dirs   []string
	delay  time.Duration
	minRun time.Duration
}

func newWatcher(t, m time.Duration) (w *watcher, err error) {
	w = &watcher{delay: t, minRun: m}
	w.Watcher, err = fsnotify.NewWatcher()
	return
}

type status struct {
	t0 time.Time
}

type runCause int

const (
	runFirst runCause = iota + 1
	runAwakened
)

func loop(w *watcher, p *prog) {

	var (
		running  bool
		hadEvent bool
		statusc  = make(chan status)
		timer    = time.NewTimer(w.delay)
	)
	timer.Stop() // need to start it with prog running or for an event

	startRunning := func(why runCause) {
		running = true

		t0 := time.Now()
		timer.Reset(w.delay)

		var causeText string
		switch why {
		case runAwakened:
			causeText = "awakened"
		case runFirst:
			causeText = "started"
		}

		logfBlue("[forever %s %s]", causeText, timef(t0))

		go func() {
			watchdebug("run->")

			pst, err := p.run()
			switch {
			case pst == nil:
				logfRed("%v", err)
			default:
				t := time.Now()
				if err != nil {
					logfRed("`%v` failed: %v", p, err)
				}
				logfBlue("[%s]", pstatef(pst, t.Sub(t0)))
			}
			statusc <- status{t0}
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
			if ev.Op&(fsnotify.Write|fsnotify.Create|fsnotify.Rename) == 0 {
				continue
			}

			watchdebug("event %v, r=%v, e=%v", ev, running, hadEvent)

			if running {
				continue
				// todo: run again with updates after a current run
			}
			hadEvent = true
			timer.Reset(w.delay)

		case st := <-statusc:
			dt := time.Now().Sub(st.t0)

			watchdebug("->st, r=%v, e=%v, st=%v, dt=%s", running, hadEvent, st, dt)

			if dt < w.minRun {
				time.AfterFunc(w.minRun-dt, func() {
					statusc <- st
				})
				continue
			}

			running = false

		case <-timer.C:
			watchdebug("timer, r=%v, e=%v", running, hadEvent)

			if hadEvent {
				hadEvent = false
				if !running {
					startRunning(runAwakened)
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

func (s status) String() string {
	return fmt.Sprintf("{t0=%s}", timef_ns(s.t0))
}

func timef(t time.Time) string {
	return t.Format("15:04:05")
}

func timef_ns(t time.Time) string {
	return t.Format("15:04:05.000000000")
}

func _watchdebug(format string, a ...any) {
	debugf(fmt.Sprintf("watch at %s: ", timef_ns(time.Now()))+format, a...)
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
