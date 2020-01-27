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
	/home/rdan/protobuf/bin/protoc -I proto -I ${GOPATH}/pkg/mod/github.com/grpc-ecosystem/grpc-gateway@v1.12.2/third_party/googleapis --go_out=plugins=grpc:proto/ proto/api.proto
	/home/rdan/protobuf/bin/protoc -I proto -I ${GOPATH}/pkg/mod/github.com/grpc-ecosystem/grpc-gateway@v1.12.2/third_party/googleapis --grpc-gateway_out=logtostderr=true:proto/ proto/api.proto
	python3 -m grpc_tools.protoc -I./proto -I ${GOPATH}/pkg/mod/github.com/grpc-ecosystem/grpc-gateway@v1.12.2/third_party/googleapis --python_out=./proto --grpc_python_out=./proto ./proto/api.proto
	/bin/cp -rf proto/*.py storage/

## help: prints this help message
help:
	@echo "Usage: \n"
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' |  sed -e 's/^/ /'
