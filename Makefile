all: test

test:
	[ -f .env ] && source .env; go test -v -count=1 ./...

build:
	go build

server:
	go run main.go server

receiver:
	go run main.go receiver

.PHONY: server receiver test
