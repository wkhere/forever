go:
	go fmt
ifneq ($(OS), Windows_NT)
	go test -race
endif
	go install

.PHONY: go
