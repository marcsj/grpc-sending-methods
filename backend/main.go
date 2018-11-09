package main

import (
	"flag"
	"fmt"
	"github.com/improbable-eng/grpc-web/go/grpcweb"
	"github.com/marcsj/streaming-grpc-web-example/backend/dog"
	"github.com/marcsj/streaming-grpc-web-example/backend/services"
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
	httpsPort := 9091
	grpcPort := 50051

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
		Addr:    fmt.Sprintf(":%d", httpsPort),
		Handler: http.HandlerFunc(handler),
	}

	errChannel := make(chan error)
	go func () {
		grpclog.Infof("Starting server. https port: %v", httpsPort)
		errChannel <- httpServer.ListenAndServeTLS(*tlsCertFilePath, *tlsKeyFilePath)
	}()
	go func () {
		grpclog.Infof("Starting server. tcp port: %v", grpcPort)
		errChannel <- grpcServer.Serve(lis)
	}()

	for {
		select {
		case <-errChannel:
			log.Fatal(<-errChannel)
		}
	}
}