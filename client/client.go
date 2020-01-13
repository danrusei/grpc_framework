package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"time"

	api "github.com/Danr17/grpc_framework/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/status"
)

var (
	addr = flag.String("addr", "localhost", "The address of the server to connect to")
	port = flag.String("port", "8080", "The port to connect to")
)

func main() {

	flag.Parse()

	if flag.NArg() < 1 {
		fmt.Fprintln(os.Stderr, "missing command: getprodtypes or getprods")
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
	case "getprods":
		err = getprods(ctx, client, flag.Arg(1), flag.Arg(2))
	default:
		err = fmt.Errorf("unknown subcommand %s", cmd)
	}
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

}

func getprodtypes(ctx context.Context, client api.ProdServiceClient, vendor string) error {

	log.Printf("requesting all product types from vendor: %s", vendor)

	if vendor == "" {
		return fmt.Errorf("Vendor arg is missing, select between available cloud vendors: google, aws, oracle")
	}

	requestProdType := api.ClientRequestType{
		Vendor: vendor,
	}

	response, err := client.GetVendorProdTypes(ctx, &requestProdType)
	if err != nil {
		if errStatus, ok := status.FromError(err); ok {
			return status.Errorf(errStatus.Code(), "error while calling client.GetVendorProdTypes() method: %v ", errStatus.Message())
		}
		return fmt.Errorf("Could not get the products: %v", err)
	}

	fmt.Printf("%s cloud products type are: %s\n", vendor, response.GetProductType())

	return nil

}

func getprods(ctx context.Context, client api.ProdServiceClient, vendor string, prodType string) error {

	log.Printf("requesting all %s products from %s", prodType, vendor)

	if vendor == "" || prodType == "" {
		return fmt.Errorf("You need both, vendor and prodType args. Example command: $client oracle storage")
	}

	requestProd := api.ClientRequestProds{
		Vendor:      vendor,
		ProductType: prodType,
	}

	stream, err := client.GetVendorProds(ctx, &requestProd)
	if err != nil {
		if errStatus, ok := status.FromError(err); ok {
			return status.Errorf(errStatus.Code(), "error while calling client.GetVendorProds() method: %v ", errStatus.Message())
		}
		return fmt.Errorf("Could not get the stream of products : %v", err)
	}

	for {
		product, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			if errStatus, ok := status.FromError(err); ok {
				return status.Errorf(errStatus.Code(), "error while receiving the stream for client.GetVendorProds: %v ", errStatus.Message())
			}
			return fmt.Errorf("error while receiving the stream for client.GetVendorProds: %v", err)
		}
		fmt.Printf("Title: %s, Url: %s,  ShortUrl: %s\n", product.GetProduct().GetTitle(), product.GetProduct().GetUrl(), product.GetProduct().GetShortUrl())
	}

	return nil
}
