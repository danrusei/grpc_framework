.DEFAULT_GOAL := help

## build: build the GRPC Client, Server (Backend is in Python)
build:
	go build -o client ./client/.
	go build -o server ./server/.

## clean: cleans all the binaries
clean:
	go clean ./...

## generate: generates the Go grpc libraries 
generate:
	/home/rdan/protobuf/bin/protoc -I proto --go_out=plugins=grpc:proto/ proto/api.proto

## help: prints this help message
help:
	@echo "Usage: \n"
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' |  sed -e 's/^/ /'
