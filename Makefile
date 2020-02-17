all: test

test:
	@[ -f .env ] && source .env; go test -count=1 ./...

build:
	go build

server:
	go run main.go server

receiver:
	cd receiver && go run main.go

.PHONY: server receiver test
