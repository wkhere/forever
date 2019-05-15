go:
	go fmt
	go test -race
	go install

.PHONY: go
