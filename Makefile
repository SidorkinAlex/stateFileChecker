.PHONY: build
build:
	go build -v ./cmd/stateFileCHecker/main.go
	mkdir -p build
	$(eval NEW_VER:=$(shell cat version | cut -d '_' -f 2 ))
	mv main build/stateFileCHecker:$(NEW_VER)

.PHONY: test
test:
	go test -v -race -timeout 30s ./...

.DEFAULT_GOAL := build
