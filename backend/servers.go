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
	"strings"
)

func getRedocHandler(handler http.Handler, basePath string, specBody []byte, spec string, specExtension string) http.Handler {
	if !strings.HasSuffix(basePath, "/") {
		basePath += "/"
	}
	specHandler := middleware.Spec(fmt.Sprintf("%sdocs/%s", basePath, spec), specBody, handler)
	return middleware.Redoc(middleware.RedocOpts{
		BasePath: basePath,
		SpecURL:  path.Join(fmt.Sprintf("%sdocs/%s", basePath, spec), specExtension),
		Path:     fmt.Sprintf("/docs/%s", spec),
	}, specHandler)
}

func getOpenAPISpecBytes(name string, path string, extension string) ([]byte, error){
	path = fmt.Sprintf("%s%s/%s.%s", path, name, name, extension)
	specDoc, err := loads.Spec(path)
	if err != nil {
		return nil, err
	}
	bytes, err := json.MarshalIndent(specDoc.Spec(), "", "  ")
	if err != nil {
		return nil, err
	}
	return bytes, nil
}

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

func getOpenAPIServer(
	port int,
	basePath string,
	specPath string,
	specExtension string,
	specs ...string) (http.Server, error) {
	handler := http.NotFoundHandler()
	for _, spec := range specs {
		bytes, err := getOpenAPISpecBytes(spec, specPath, specExtension)
		if err != nil {
			return http.Server{}, err
		}
		handler = getRedocHandler(handler, basePath, bytes, spec, specExtension)
	}

	server := http.Server{
		Addr: fmt.Sprintf(":%v", port),
		Handler: handler,
	}
	return server, nil
}