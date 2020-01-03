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
	prod map[string][]string
}

func newServer() *server {

	serv := map[string][]string{
		"google": []string{"compute", "storage"},
		"aws":    []string{"compute", "storage"},
		"oracle": []string{"compute", "storage"},
	}

	return &server{prod: serv}
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

	api.RegisterProdServiceServer(srv, newServer())

	log.Printf("Serving gRPC on https://%s", addr)

	if err := (srv.Serve(lis)); err != nil {
		return fmt.Errorf("Unable to start GRPC server: %s", err)
	}

	return nil
}

//GetProds implement the GRPC server function
func (serv server) GetProds(ctx context.Context, req *api.ClientRequest) (*api.ClientResponse, error) {

	var prods string

	if vendorProds, found := serv.prod[req.GetVendor()]; found {

		for _, prod := range vendorProds {
			prods = prods + " " + prod
		}

	} else {
		return nil, fmt.Errorf("vendor unavailable")
	}

	clientResponse := api.ClientResponse{
		Products: prods,
	}

	return &clientResponse, nil
}
