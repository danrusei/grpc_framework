package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"time"

	api "github.com/Danr17/grpc_framework/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/status"
)

var (
	port = flag.Int("grpc-port", 8080, "The gRPC server port")
)

type server struct {
	prodTypes   map[string][]string
	storageAddr string
	storagePort string
}

func newServer(vendorServ map[string][]string) *server {

	return &server{
		prodTypes:   vendorServ,
		storageAddr: "localhost",
		storagePort: "6000",
	}
}

var vendorServices = map[string][]string{
	"google": []string{"compute", "storage"},
	"aws":    []string{"compute", "storage"},
	"oracle": []string{"compute", "storage"},
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

	creds, err := credentials.NewServerTLSFromFile("../cert/service.pem", "../cert/service.key")
	if err != nil {
		return fmt.Errorf("could not process the credentials: %v", err)
	}

	// Create an array of gRPC options with the credentials
	opts := []grpc.ServerOption{grpc.Creds(creds)}

	srv := grpc.NewServer(opts...)

	api.RegisterProdServiceServer(srv, newServer(vendorServices))

	log.Printf("Serving gRPC on https://%s", addr)

	if err := (srv.Serve(lis)); err != nil {
		return fmt.Errorf("Unable to start GRPC server: %s", err)
	}

	return nil
}

//GetVendorProdTypes implement the GRPC server function
func (serv *server) GetVendorProdTypes(ctx context.Context, req *api.ClientRequestType) (*api.ClientResponseType, error) {

	log.Printf("have received a request for -> %s <- as vendor", req.GetVendor())

	var prodTypes string

	if vendorProdTypes, found := serv.prodTypes[req.GetVendor()]; found {

		for _, prodType := range vendorProdTypes {
			prodTypes = prodTypes + " " + prodType
		}

	} else {
		return nil, status.Errorf(codes.InvalidArgument, "wrong vendor, select between google, aws, oracle")
	}

	// some heavy processing **increase it ** if you want to test out DeadlineExceeded
	time.Sleep(2 * time.Second)

	if ctx.Err() == context.DeadlineExceeded {
		log.Printf("dealine has exceeded, stoping server side operation")
		return nil, status.Error(codes.DeadlineExceeded, "dealine has exceeded, stoping server side operation")
	}

	if ctx.Err() == context.Canceled {
		log.Print("the user has canceled the request, stoping server side operation")
		return nil, status.Error(codes.Canceled, "the user has canceled the request, stoping server side operation")
	}

	clientResponse := api.ClientResponseType{
		ProductType: prodTypes,
	}

	log.Printf("the response is sent to client: %s", prodTypes)

	return &clientResponse, nil
}

func (serv *server) GetVendorProds(req *api.ClientRequestProds, stream api.ProdService_GetVendorProdsServer) error {

	log.Printf("have received a request for -> %s <- product type from -> %s <- vendor", req.GetProductType(), req.GetVendor())

	conn, err := grpc.Dial(net.JoinHostPort(serv.storageAddr, serv.storagePort), grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to dial server:, %s", err)

	}
	defer conn.Close()

	ctx := context.Background()

	client := api.NewStorageServiceClient(conn)
	response, err := client.GetProdsDetail(ctx, &api.StorageRequest{
		Vendor:      req.GetVendor(),
		ProductType: req.GetProductType(),
	})
	if err != nil {
		return status.Errorf(codes.Internal, "error while calling client.GetProdsDetail() method: %v ", err)
	}

	for _, prod := range response.ProdDetail {
		if err := stream.Send(&api.ClientResponseProds{
			Product: &api.ProdsPrep{
				Title:    prod.GetTitle(),
				Url:      prod.GetUrl(),
				ShortUrl: "None- TBD",
			},
		}); err != nil {
			return status.Error(codes.Internal, "not able to send the response")
		}
	}

	return nil
}
