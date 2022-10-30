.PHONY: build

build:
	go build -o exe/snake ./cmd/snake
	go build -o exe/xade ./cmd/xade
