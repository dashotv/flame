
all: test

test: build
	test ! -z "$$FLAME_URL"
	go test

build:
	go build
