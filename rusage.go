package main

import "os"

type memStats struct {
	maxRss         int64
	minFlt, majFlt int64
}

type rusageExtrasGetter interface {
	getMemStats(*os.ProcessState) (memStats, bool)
}

var rusageExtras rusageExtrasGetter = noRusageExtras{}

type noRusageExtras struct{}

func (noRusageExtras) getMemStats(*os.ProcessState) (memStats, bool) {
	return memStats{}, false
}
