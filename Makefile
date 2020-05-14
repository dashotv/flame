all: test

test:
	[ -f .env ] && source .env; go test -v -count=1 ./...

generate:
	golem generate

build: generate
	go build

server: generate
	go run main.go server

receiver:
	go run main.go receiver

.PHONY: server receiver test
