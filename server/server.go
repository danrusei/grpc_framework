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
	port = flag.Int("grpc-port", 8080, "The gRPC server port")
)

type server struct {
	prodTypes map[string][]string
}

func newServer(vendorServ map[string][]string) *server {

	return &server{prodTypes: vendorServ}
}

var vendorServices = map[string][]string{
	"google": []string{"google_compute", "google_storage"},
	"aws":    []string{"aws_compute", "aws_storage"},
	"oracle": []string{"oracle_compute", "oracle_storage"},
}

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

	api.RegisterProdServiceServer(srv, newServer(vendorServices))

	log.Printf("Serving gRPC on https://%s", addr)

	if err := (srv.Serve(lis)); err != nil {
		return fmt.Errorf("Unable to start GRPC server: %s", err)
	}

	return nil
}

//GetVendorProdTypes implement the GRPC server function
func (serv server) GetVendorProdTypes(ctx context.Context, req *api.ClientRequestType) (*api.ClientResponseType, error) {

	log.Printf("have received a request for -> %s <- as vendor", req.GetVendor())

	var prodTypes string

	if vendorProdTypes, found := serv.prodTypes[req.GetVendor()]; found {

		for _, prodType := range vendorProdTypes {
			prodTypes = prodTypes + " " + prodType
		}

	} else {
		return nil, fmt.Errorf("wrong vendor selection")
	}

	clientResponse := api.ClientResponseType{
		ProductType: prodTypes,
	}

	return &clientResponse, nil
}
