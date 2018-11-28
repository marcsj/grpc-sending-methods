# gRPC Sending Methods

Examples of streaming and otherwise sending gRPC from a go backend to clients using a dog tracker. 
Data is sent from the backend via these methods:
- **gRPC**
- **grpc-gateway**(JSON over http) 
- **grpc-gateway with websockets**(JSON over ws), 
- **grpc-web**(http2)


## Server:

1. Set up `$GOPATH` and add `$GOPATH/bin` to your `$PATH`
2. Fetch the following dependencies:
```
go get -u github.com/golang/protobuf/protoc-gen-go
go get -u github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway
go get -u github.com/grpc-ecosystem/grpc-gateway/protoc-gen-swagger
go get -u github.com/golang/protobuf/protoc-gen-go
```
3. `cd proto`
4. Run `./generate.sh`
5. For Go 1.11 vendoring, set envionmental variable `GO111MODULE=on`.
6. Run `go mod vendor` to update vendor folder.
7. `go build main.go` or `go run main.go` in backend directory.


## Clients:

### Web-WS
Javascript web client using json over websockets for streaming
1. `npm i`
2. `npm start`
3. Navigate to [http://localhost:8080](http://localhost:8080)

### Web-GRPC
Javascript web client using grpc over grpc-web
1. `npm i`
2. `npm start`
3. Navigate to [https://localhost:8080](http://localhost:8080)

### Use:

Look at the backend logs to see the initialization of each dog location ID. 
The IDs listed can be entered into the client and used to receive dogs that are currently in those locations.


## Certs:

Certs are needed for the proxy and for the browser that accesses it via grpc-web.

- Go to [github.com/kingkool68/generate-ssl-certs-for-local-development](https://github.com/kingkool68/generate-ssl-certs-for-local-development) and follow their instructions for generating certs.
- Certs should be placed under `keys/`.
