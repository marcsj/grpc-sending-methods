package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/marcsj/grpc-sending-methods/backend/dog"
	"github.com/marcsj/grpc-sending-methods/backend/services"
	"github.com/marcsj/grpc-sending-methods/backend/store"
	"github.com/tmc/grpc-websocket-proxy/wsproxy"
	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"
	"log"
	"net"
	"net/http"
	"os"
)

func init() {
	grpcLog := grpclog.NewLoggerV2(os.Stdout, os.Stderr, os.Stderr)
	grpclog.SetLoggerV2(grpcLog)
}

var grpcPort = flag.Int(
	"grpc_port",
	50051,
	"port for gRPC calls")
var grpcwPort = flag.Int(
	"grpcw_port",
	9091,
	"port for grpc-web")
var gatewayPort = flag.Int(
	"gateway_port",
	8081,
	"port for grpc-gateway requests")
var openAPIPort = flag.Int(
	"openapi_port",
	8082,
	"port for open api")
var numberDogs = flag.Int(
	"number_dogs",
	1000,
	"mumber of dogs to generate")
var numberDaycares = flag.Int(
	"number_daycares",
	4,
	"number of daycares to generate")

func main() {
	errChannel := make(chan error)

	// setup for gRPC server
	lis, err := net.Listen("tcp", fmt.Sprintf(":%v", *grpcPort))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()

	// setup services
	dogStore := store.NewDogStore(*numberDaycares, *numberDogs)
	dogTrackServer := services.NewDogTrackServer(dogStore)
	dog.RegisterDogTrackServer(grpcServer, dogTrackServer)

	// setup for proxy for grpc-web
	grpcWebServer := getGRPCWebServer(grpcServer, *grpcwPort)

	// setup for openAPI server
	openAPIServer, err := getOpenAPIServer(
		"/", "", "swagger.json", "dog")
	if err != nil {
		log.Fatal(err)
	}

	// setup for gRPC-gateway
	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithInsecure()}

	// running gRPC server
	go func () {
		grpclog.Infof("Starting gRPC server. tcp port: %v", *grpcPort)
		errChannel <- grpcServer.Serve(lis)
	}()

	// running proxy for grpc-web
	go func () {
		grpclog.Infof("Starting grpc-web proxy server. http port: %v", *grpcwPort)
		errChannel <- grpcWebServer.ListenAndServe()
	}()

	// running openAPI server to show api info
	go func () {
		grpclog.Infof("Starting OpenAPI server. http port: %v", *openAPIPort)
		errChannel <- openAPIServer.ListenAndServe()
	}()

	// running gRPC-gateway
	logger := log.New(os.Stdout, "http: ", log.LstdFlags)
	go func () {
		errChannel <- dog.RegisterDogTrackHandlerFromEndpoint(
			context.Background(), mux, fmt.Sprintf("localhost:%v", *grpcPort), opts)
		grpclog.Infof("Starting gRPC-gateway server. http port: %v", *gatewayPort)
		grpcGateway := http.Server{
			Addr: fmt.Sprintf(":%v", *gatewayPort),
			Handler: wsproxy.WebsocketProxy(mux, wsproxy.WithMethodParamOverride("method")),
			ErrorLog: logger,
		}
		errChannel <- grpcGateway.ListenAndServe()
	}()

	for {
		select {
		case err := <-errChannel:
			if err != nil {
				log.Fatal(err)
			}
		}
	}
}