# Streaming gRPC-web Example

An example of grpc-web using a streaming dog tracker. This is a simple example of streaming grpc to grpc-web, and does not use an actual datastore apart from the list of dogs created in memory.

## Building:

1. Set up `$GOPATH` and add `$GOPATH/bin` to your `$PATH`
2. Run `go get -u github.com/golang/protobuf/protoc-gen-go`
3. `cd proto`
4. Run `./generate.sh`
5. For Go 1.11 vendoring, set envionmental variable `GO111MODULE=on`.
6. Run `go mod vendor` to update vendor folder.
7. `go build` or `go run`.

## Certs:

Certs are needed for the proxy and for the browser that accesses it via grpc-web.

- Go to [github.com/kingkool68/generate-ssl-certs-for-local-development](https://github.com/kingkool68/generate-ssl-certs-for-local-development) and follow their instructions for generating certs.
- Certs should be placed under `keys/`.
