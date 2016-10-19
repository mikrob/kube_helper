BINARY=kube_helper

test:
	go test ./...

all:
	go build -o ${BINARY} main.go
	sudo cp ${BINARY} /usr/local/bin/
