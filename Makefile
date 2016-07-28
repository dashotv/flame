
all: test

test: build
	go test

build: deps
	go build

deps:
	go get github.com/Masterminds/glide
	glide install
