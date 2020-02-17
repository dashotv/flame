FLAME_URL = "http://admin:utorrent@utorrent.dasho.net/gui"

all: test

test:
	go test ./...

build:
	go build

server:
	go run main.go server

receiver:
	cd receiver && go run main.go

.PHONY: server receiver test
