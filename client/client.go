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
	addr   = flag.String("addr", "localhost", "The address of the server to connect to")
	port   = flag.String("port", "8080", "The port to connect to")
	vendor = flag.String("vendor", "", "Select a vendor")
	pType  = flag.String("ptype", "", "Select a product type")
)

func main() {

	flag.Parse()

	if flag.NArg() < 1 {
		fmt.Fprintln(os.Stderr, "missing command: getprodtypes")
		os.Exit(1)
	}

	conn, err := grpc.Dial(net.JoinHostPort(*addr, *port), grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to dial server:, %s", err)

	}
	defer conn.Close()

	client := api.NewProdServiceClient(conn)

	ctx := context.Background()

	switch cmd := flag.Arg(0); cmd {
	case "getprodtypes":
		err = getprodtypes(ctx, client, *vendor)
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
		return fmt.Errorf("Vendor flag is missing, select between available cloud vendors: google, aws, oracle")
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
