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
	go get golang.org/x/tools/cmd/goimports
	go get github.com/dashotv/golem

.PHONY: server receiver test deps docker docker-run
