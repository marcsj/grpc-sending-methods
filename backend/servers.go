package main

import (
	"encoding/json"
	"fmt"
	"github.com/go-openapi/loads"
	"github.com/go-openapi/runtime/middleware"
	"github.com/improbable-eng/grpc-web/go/grpcweb"
	"google.golang.org/grpc"
	"net/http"
	"path"
)

func getGRPCWebServer(grpcServer *grpc.Server, port int) (http.Server) {
	wrappedServer := grpcweb.WrapServer(grpcServer)
	handler := func(resp http.ResponseWriter, req *http.Request) {
		wrappedServer.ServeHTTP(resp, req)
	}
	return http.Server{
		Addr:    fmt.Sprintf(":%v", port),
		Handler: http.HandlerFunc(handler),
	}
}

func getOpenAPIServer(basePath string, specs ...string) (http.Server, error){
	newestHandler := http.NotFoundHandler()
	for i, s := range specs {
		specDoc, err := loads.Spec(
			fmt.Sprintf("%s/%s.swagger.json", s, s))
		if err != nil {
			return http.Server{}, err
		}
		b, err := json.MarshalIndent(specDoc.Spec(), "", "  ")
		if err != nil {
			return http.Server{}, err
		}

		if i == 0 {
			specHandler := middleware.Spec(fmt.Sprintf("/%s", s), b, nil)
			newestHandler = middleware.Redoc(middleware.RedocOpts{
				BasePath: basePath,
				SpecURL:  path.Join(fmt.Sprintf("/%s", s), "swagger.json"),
				Path:     fmt.Sprintf("docs/%s", s),
			}, specHandler)
		} else {
			specHandler := middleware.Spec(fmt.Sprintf("/%s", s), b, newestHandler)
			newestHandler = middleware.Redoc(middleware.RedocOpts{
				BasePath: basePath,
				SpecURL:  path.Join(fmt.Sprintf("/%s", s), "swagger.json"),
				Path:     fmt.Sprintf("docs/%s", s),
			}, specHandler)
		}
	}

	openAPIServer := http.Server{
		Addr: fmt.Sprintf(":%v", *openAPIPort),
		Handler: newestHandler,
	}
	return openAPIServer, nil
}