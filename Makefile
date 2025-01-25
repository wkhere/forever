go:
ifneq ($(OS), Windows_NT)
	go test -race
endif
	go build

install:
	go install

.PHONY: go install
