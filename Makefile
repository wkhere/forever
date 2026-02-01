go:
ifneq ($(OS), Windows_NT)
	go test -race .
endif
	go build

static:
	CGO_ENABLED=0 go build -o forever.static

install: go
	go install -ldflags=-s

.PHONY: go install static
