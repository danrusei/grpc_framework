package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"time"

	api "github.com/Danr17/grpc_framework/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

var (
	addr = flag.String("addr", "localhost", "The address of the server to connect to")
	port = flag.String("port", "8080", "The port to connect to")
)

func main() {

	flag.Parse()

	if flag.NArg() < 1 {
		fmt.Fprintln(os.Stderr, "missing command: getprodtypes")
		os.Exit(1)
	}

	creds, err := credentials.NewClientTLSFromFile("../cert/service.pem", "")
	if err != nil {
		log.Fatalf("could not process the credentials: %v", err)
	}

	conn, err := grpc.Dial(net.JoinHostPort(*addr, *port), grpc.WithTransportCredentials(creds))
	if err != nil {
		log.Fatalf("Failed to dial server:, %s", err)

	}
	defer conn.Close()

	client := api.NewProdServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
	defer cancel()

	switch cmd := flag.Arg(0); cmd {
	case "getprodtypes":
		err = getprodtypes(ctx, client, flag.Arg(1))
	default:
		err = fmt.Errorf("unknown subcommand %s", cmd)
	}
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

}

func getprodtypes(ctx context.Context, client api.ProdServiceClient, vendor string) error {

	if vendor == "" {
		return fmt.Errorf("Vendor arg is missing, select between available cloud vendors: google, aws, oracle")
	}

	requestProd := api.ClientRequestType{
		Vendor: vendor,
	}

	response, err := client.GetVendorProdTypes(ctx, &requestProd)
	if err != nil {
		return fmt.Errorf("Could not get the products: %v", err)
	}

	fmt.Printf("%s cloud products type are: %s\n", vendor, response.GetProductType())

	return nil

}
