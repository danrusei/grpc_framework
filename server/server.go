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
	storageAddr string
	storagePort string
}

func newServer() *server {
	serv := server{
		storageAddr: "localhost",
		storagePort: "6000",
	}
	return &serv
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

	conn, err := grpc.DialContext(ctx, net.JoinHostPort(serv.storageAddr, serv.storagePort), grpc.WithInsecure())
	if err != nil {
		return nil, fmt.Errorf("Failed to dial server:, %s", err)
	}
	defer conn.Close()

	storage := api.NewBackendServiceClient(conn)

	itemsRequest := api.ApiRequest{
		Vendor:   req.Vendor,
		ProdType: req.ProdType,
	}
	itemsResponse, err := storage.RetrieveItems(ctx, &itemsRequest)
	if err != nil {
		return nil, fmt.Errorf("wrong RPC request: %v", err)
	}

	var prods string

	for _, prod := range itemsResponse.Prods {
		prods = prods + prod
	}

	clientResponse := api.ClientResponse{
		Products: prods,
	}

	return &clientResponse, nil
}
