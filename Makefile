
all: test

build:
	go build client.go response.go torrent.go

test:
	go test

update:
	gvt update -all
