
all: test

build:
	go build client.go response.go torrent.go

test:
	go test

deps:
	gvt restore

update:
	gvt update -all
