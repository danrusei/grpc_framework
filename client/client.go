package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"os"

	api "github.com/Danr17/grpc_framework/proto"
	"google.golang.org/grpc"
)

var (
	addr = flag.String("addr", "localhost", "The address of the server to connect to")
	port = flag.String("port", "8080", "The port to connect to")
)

func main() {

	flag.Parse()

	if flag.NArg() < 1 {
		fmt.Fprintln(os.Stderr, "missing vendor: google, aws, oracle")
		os.Exit(1)
	}

	if err := run(*addr, *port, flag.Arg(0)); err != nil {
		log.Fatalln(err)
	}

}

func run(addr, port string, vendor string) error {
	ctx := context.Background()

	conn, err := grpc.DialContext(ctx, net.JoinHostPort(addr, port), grpc.WithInsecure())
	if err != nil {
		return fmt.Errorf("Failed to dial server:, %s", err)
	}
	defer conn.Close()

	requestProd := api.ClientRequestType{
		Vendor: vendor,
	}

	client := api.NewProdServiceClient(conn)
	response, err := client.GetVendorProdTypes(ctx, &requestProd)
	if err != nil {
		return fmt.Errorf("Could not get the products: %v", err)
	}

	fmt.Printf("%s cloud products type are: %s\n", vendor, response.GetProductType())

	return nil

}
