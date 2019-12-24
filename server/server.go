package main

import (
	"flag"
	"fmt"
	"google.golang.org/grpc"
	"log"
	"net"
)

var (
	port = flag.Int("grpc-port", 8080, "The gRPC server port")
)

func main() {

	flag.Parse()
	addr := fmt.Sprintf("localhost:%d", *port)
	if err := run(addr); err != nil {
		log.Fatalf("could not start the server: %s", err)
	}
}

func run(addr string) error {
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("could not listen on the port %s: %s", addr, err)
	}

	srv := grpc.NewServer()

	if err := (srv.Serve(lis)); err != nil {
		return fmt.Errorf("Unable to start GRPC server: %s", err)
	}

	return nil
}
