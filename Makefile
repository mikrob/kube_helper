BINARY=kube_helper

PHONY: all

test:
	go test  -v ./...

get:
	go get

all:
	go build -o ${BINARY} main.go
	sudo cp ${BINARY} /usr/local/bin/
