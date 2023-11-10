NAME := shawncatz/flame
PORT := 9001

all: test

test:
	[ -f .env ] && source .env; go test -v ./qbt ./nzbget

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
	go install github.com/codegangsta/gin@latest

.PHONY: server receiver test deps docker docker-run
