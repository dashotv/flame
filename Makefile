
all: test

test: build
	test ! -z "$$FLAME_URL"
	go test

build: deps
	go build

deps:
	go get github.com/Masterminds/glide
	glide install
