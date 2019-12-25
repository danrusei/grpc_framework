.DEFAULT_GOAL := help

## build: build the GRPC Client, Server and Backend
build:
	go build -o client ./client/.
	go build -o server ./server/.
	go build -o backend ./backend/.

## clean: cleans all the binaries
clean:
	go clean ./...

## generate: generates the Go and Python grpc libraries 
generate:
	/home/rdan/protobuf/bin/protoc -I proto --go_out=plugins=grpc:proto/ proto/api.proto
	python3 -m grpc_tools.protoc -I./proto --python_out=./proto --grpc_python_out=./proto ./proto/api.proto

## help: prints this help message
help:
	@echo "Usage: \n"
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' |  sed -e 's/^/ /'
