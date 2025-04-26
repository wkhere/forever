go:
ifneq ($(OS), Windows_NT)
	go test
endif
	go build

static:
	CGO_ENABLED=0 go build -o forever.static

install:
	go install

.PHONY: go install static
