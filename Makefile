BINARY=kube_helper

all:
	go build -o ${BINARY} main.go
	sudo cp ${BINARY} /usr/local/bin/

