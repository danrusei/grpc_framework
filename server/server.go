package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/Danr17/grpc_framework/middleware/grpcklog"
	"github.com/Danr17/grpc_framework/middleware/grpcopentelemetry"
	api "github.com/Danr17/grpc_framework/proto"
	"github.com/go-logr/logr"
	"k8s.io/klog/klogr"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/status"

	"github.com/google/uuid"
)

var (
	portGRPC = flag.Int("grpc-port", 8080, "The gRPC server port")
	portREST = flag.Int("rest-port", 8081, "The REST server port")
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
	logger := klogr.New()
	grpcopentelemetry.Init()

	grpcAddr := fmt.Sprintf("localhost:%d", *portGRPC)
	go func() {
		if err := runGRPCServer(logger, grpcAddr); err != nil {
			log.Fatalf("could not start the server: %s", err)
		}
	}()

	restAddr := fmt.Sprintf("localhost:%d", *portREST)
	go func() {
		if err := runRESTServer(restAddr, grpcAddr); err != nil {
			log.Fatalf("could not start the server: %s", err)
		}
	}()

	// infinite loop
	log.Printf("Entering infinite loop")
	select {}
}

func runRESTServer(restAddr, grpcAdd string) error {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	mux := runtime.NewServeMux()
	creds, err := credentials.NewClientTLSFromFile("../cert/service.pem", "")
	if err != nil {
		return fmt.Errorf("could not load TLS certificate: %s", err)
	}
	// Setup the client gRPC options
	opts := []grpc.DialOption{grpc.WithTransportCredentials(creds)}
	// Register ping
	err = api.RegisterProdServiceHandlerFromEndpoint(ctx, mux, grpcAdd, opts)
	if err != nil {
		return fmt.Errorf("could not register service ProdService: %s", err)
	}
	log.Printf("starting HTTP/1.1 REST server on %s", restAddr)
	http.ListenAndServe(restAddr, mux)
	return nil
}

func runGRPCServer(logger logr.Logger, addr string) error {
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("could not listen on the port %s: %s", addr, err)
	}

	creds, err := credentials.NewServerTLSFromFile("../cert/service.pem", "../cert/service.key")
	if err != nil {
		return fmt.Errorf("could not process the credentials: %v", err)
	}

	// Shared options for the logger, with a custom duration to log field function.
	optsLog := []grpcklog.Option{
		grpcklog.WithDurationField(func(duration time.Duration) (key string, value interface{}) {
			return "grpc.time_ns", duration.Nanoseconds()
		}),
	}

	srv := grpc.NewServer(
		grpc_middleware.WithUnaryServerChain(
			grpc_ctxtags.UnaryServerInterceptor(grpc_ctxtags.WithFieldExtractor(grpc_ctxtags.CodeGenRequestFieldExtractor)),
			grpcklog.UnaryServerInterceptor(logger, optsLog...),
			grpcopentelemetry.UnaryServerInterceptor,
		),
		grpc.Creds(creds),
		grpc_middleware.WithStreamServerChain(
			grpc_ctxtags.StreamServerInterceptor(grpc_ctxtags.WithFieldExtractor(grpc_ctxtags.CodeGenRequestFieldExtractor)),
			grpcklog.StreamServerInterceptor(logger, optsLog...),
		),
	)

	api.RegisterProdServiceServer(srv, newServer(vendorServices))

	log.Printf("Serving gRPC on https://%s", addr)

	if err := (srv.Serve(lis)); err != nil {
		return fmt.Errorf("Unable to start GRPC server: %s", err)
	}

	return nil
}

//GetVendorProdTypes implement the GRPC server function
func (serv *server) GetVendorProdTypes(ctx context.Context, req *api.ClientRequestType) (*api.ClientResponseType, error) {
	//log.Printf("have received a request for -> %s <- as vendor", req.GetVendor())
	//let's assume we are able to identify the calling Customer, fake it with random numbers
	addCustomerToctx(ctx)

	var prodTypes string
	if vendorProdTypes, found := serv.prodTypes[req.GetVendor()]; found {

		for _, prodType := range vendorProdTypes {
			prodTypes = prodTypes + " " + prodType
		}
	} else {
		return nil, status.Errorf(codes.InvalidArgument, "wrong vendor, select between google, aws, oracle")
	}

	//to simulate heavy processing **increase it ** -- to test out DeadlineExceeded
	//time.Sleep(1 * time.Second)

	if ctx.Err() == context.DeadlineExceeded {
		//	log.Printf("dealine has exceeded, stoping server side operation")
		return nil, status.Error(codes.DeadlineExceeded, "dealine has exceeded, stoping server side operation")
	}
	if ctx.Err() == context.Canceled {
		//	log.Print("the user has canceled the request, stoping server side operation")
		return nil, status.Error(codes.Canceled, "the user has canceled the request, stoping server side operation")
	}

	clientResponse := api.ClientResponseType{
		ProductType: prodTypes,
	}

	return &clientResponse, nil
}

func (serv *server) GetVendorProds(req *api.ClientRequestProds, stream api.ProdService_GetVendorProdsServer) error {

	log.Printf("have received a request for -> %s <- product type from -> %s <- vendor", req.GetProductType(), req.GetVendor())
	//let's assume we are able to identify the calling Customer, fake it with random numbers
	addCustomerToctx(stream.Context())

	conn, err := grpc.Dial(net.JoinHostPort(serv.storageAddr, serv.storagePort), grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to dial server:, %s", err)

	}
	defer conn.Close()

	ctx := stream.Context()
	client := api.NewStorageServiceClient(conn)
	response, err := client.GetProdsDetail(ctx, &api.StorageRequest{
		Vendor:      req.GetVendor(),
		ProductType: req.GetProductType(),
	})
	if err != nil {
		if errStatus, ok := status.FromError(err); ok {
			//	log.Printf("error while calling client.GetProdsDetail() method: %v ", errStatus.Message())
			return status.Errorf(errStatus.Code(), "error while calling client.GetProdsDetail() method: %v ", errStatus.Message())
		}
		//	log.Printf("error while calling client.GetProdsDetail() method: %v ", err)
		return status.Errorf(codes.Internal, "error while calling client.GetProdsDetail() method: %v ", err)
	}

	for _, prod := range response.ProdDetail {

		id := uuid.Must(uuid.NewRandom()).String()

		if err := stream.Send(&api.ClientResponseProds{
			Product: &api.ProdsPrep{
				Title:    prod.GetTitle(),
				Url:      prod.GetUrl(),
				ShortUrl: "https://made-up-url.com/" + id[:6],
			},
		}); err != nil {
			return status.Error(codes.Internal, "not able to send the response")
		}

		// to simulate heavy processing **increase it ** -- to test out DeadlineExceeded
		//time.Sleep(1 * time.Second)

		if ctx.Err() == context.DeadlineExceeded {
			//	log.Printf("dealine has exceeded, stoping server side operation")
			return status.Error(codes.DeadlineExceeded, "dealine has exceeded, stoping server side operation")
		}
		if ctx.Err() == context.Canceled {
			//	log.Print("the user has canceled the request, stoping server side operation")
			return status.Error(codes.Canceled, "the user has canceled the request, stoping server side operation")
		}
	}

	return nil
}

func addCustomerToctx(ctx context.Context) {
	clientID := uuid.Must(uuid.NewRandom()).String()
	grpcklog.AddFields(ctx, map[string]interface{}{"Name": "Customer-0367" + clientID[:4]})
}
