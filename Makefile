NAME := flame
PORT := 9001

all: test

test: generate
	[ -f .env ] && source .env; go test -v -count=1 ./qbt ./nzbget

generate:
	golem generate

build: generate
	go build

install: build
	go install

server:
	go run main.go server

receiver:
	go run main.go receiver

docker:
	docker build -t $(NAME) .

docker-run:
	docker run -d --rm --name $(NAME) -p $(PORT):$(PORT) $(NAME)

deps:
	go install golang.org/x/tools/cmd/goimports@latest
	go install github.com/dashotv/golem@latest

.PHONY: server receiver test deps docker docker-run
