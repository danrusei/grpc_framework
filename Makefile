.DEFAULT_GOAL := help

## build: build the GRPC Client, Server and Backend
build:
	go build -o client ./client/.
	go build -o server ./server/.
	go build -o backend ./backend/.

## clean: cleans all the binaries
clean:
	go clean ./...

## help: prints this help message
help:
	@echo "Usage: \n"
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' |  sed -e 's/^/ /'
