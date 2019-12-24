package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"
)

var (
	addr = flag.String("addr", "localhost", "The address of the server to connect to")
	port = flag.String("port", "10000", "The port to connect to")
)

func main() {

	flag.Parse()

	if err := run(*addr, *port); err != nil {
		log.Fatalf("could not start the client: %s", err)
	}

}

func run(addr, port string) error {
	ctx := context.Background()

	conn, err := grpc.DialContext(ctx, net.JoinHostPort(addr, port), grpc.WithInsecure())
	if err != nil {
		return fmt.Errorf("Failed to dial server:, %s", err)
	}
	defer conn.Close()

	return nil

}
