package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/improbable-eng/grpc-web/go/grpcweb"
	"github.com/marcsj/streaming-grpc-web-example/backend/dog"
	"github.com/marcsj/streaming-grpc-web-example/backend/services"
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

var tlsCertFilePath = flag.String(
	"tls_cert_file",
	"../keys/local.dev.crt",
	"Path to the CRT/PEM file.")
var tlsKeyFilePath = flag.String(
	"tls_key_file",
	"../keys/local.dev.key",
	"Path to the private key file.")

func main() {
	// initializes our sample setup
	services.InitializeDogs()
	grpcPort := 50051
	httpsPort := 9091
	gatewayPort := 8080

	// setup for gRPC server
	lis, err := net.Listen("tcp", fmt.Sprintf(":%v", grpcPort))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	dog.RegisterDogTrackServer(grpcServer, services.DogTrackServer{})

	// setup for proxy for grpc-web
	wrappedServer := grpcweb.WrapServer(grpcServer)
	handler := func(resp http.ResponseWriter, req *http.Request) {
		wrappedServer.ServeHTTP(resp, req)
	}
	httpServer := http.Server{
		Addr:    fmt.Sprintf(":%v", httpsPort),
		Handler: http.HandlerFunc(handler),
	}
	logger := log.New(os.Stdout, "http: ", log.LstdFlags)

	// setup for gRPC-gateway
	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithInsecure()}

	errChannel := make(chan error)

	// running gRPC server
	go func () {
		grpclog.Infof("Starting gRPC server. tcp port: %v", grpcPort)
		errChannel <- grpcServer.Serve(lis)
	}()

	// running proxy for grpc-web
	go func () {
		grpclog.Infof("Starting grpc-web proxy server. https port: %v", httpsPort)
		errChannel <- httpServer.ListenAndServeTLS(*tlsCertFilePath, *tlsKeyFilePath)
	}()

	// running gRPC-gateway
	go func () {
		errChannel <- dog.RegisterDogTrackHandlerFromEndpoint(
			context.Background(), mux, fmt.Sprintf("localhost:%v", grpcPort), opts)
		grpclog.Infof("Starting gRPC-gateway server. https port: %v", gatewayPort)
		grpcGateway := http.Server{
			Addr: fmt.Sprintf(":%v", gatewayPort),
			Handler: wsproxy.WebsocketProxy(mux),
			ErrorLog: logger,
		}
		errChannel <- grpcGateway.ListenAndServe()
	}()

	for {
		select {
		case <-errChannel:
			log.Fatal(<-errChannel)
		}
	}
}