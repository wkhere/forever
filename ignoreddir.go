package main

import "path/filepath"

var ignoredDirs = []string{
	".git",
	".stfolder",
	".stversions",
	"vendor",
	"__pycache__",
	".mypy_cache",
	"ebin",
	"deps",
	"_build",
	"classes",
}

func isIgnoredDir(path string) bool {
	_, last := filepath.Split(path)
	for _, d := range ignoredDirs {
		if d == last {
			return true
		}
	}
	return false
}
