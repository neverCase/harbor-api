.PHONY: mod

mod:
	go mod download
	go mod tidy