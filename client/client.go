package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"

	api "github.com/Danr17/grpc_framework/proto"
	"google.golang.org/grpc"
)

var (
	addr    = flag.String("addr", "localhost", "The address of the server to connect to")
	port    = flag.String("port", "10000", "The port to connect to")
	vendor  = flag.String("vendor", "google", "Select a cloud Vendor")
	product = flag.String("product", "", "Select a produc from available options: compute, database")
)

func main() {

	flag.Parse()

	if err := run(*addr, *port, *vendor, *product); err != nil {
		log.Fatalf("could not start the client: %s", err)
	}

}

func run(addr, port string, vendor string, product string) error {
	ctx := context.Background()

	conn, err := grpc.DialContext(ctx, net.JoinHostPort(addr, port), grpc.WithInsecure())
	if err != nil {
		return fmt.Errorf("Failed to dial server:, %s", err)
	}
	defer conn.Close()

	requestProd := api.ClientRequest{
		Vendor:   vendor,
		ProdType: product,
	}

	client := api.NewProdServiceClient(conn)
	response, err := client.GetProds(ctx, &requestProd)
	if err != nil {
		return nil
	}

	fmt.Println(response)

	return nil

}
