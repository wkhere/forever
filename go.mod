module github.com/wkhere/forever

go 1.12

require (
	github.com/fsnotify/fsnotify v1.4.9
	github.com/spf13/pflag v1.0.6
)

// fix vuln:
require golang.org/x/sys v0.0.0-20220412211240-33da011f77ad // indirect
