# gRPC Sending Methods

Examples of streaming and otherwise sending gRPC from a go backend to clients using a dog tracker. 
Data is sent from the backend via these methods:
- **[gRPC](https://github.com/grpc/grpc-go)**
- **[gRPC-gateway](https://github.com/grpc-ecosystem/grpc-gateway)**(JSON over http)
- **gRPC-gateway with websockets**(JSON over ws), 
- **[grpc-web](https://github.com/improbable-eng/grpc-web)**

These sending methods enable you to generate interfaces and models with gRPC once,
 and to be called in multiple different ways. The most important way this helps 
 you leverage code is by streaming to a web client.
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
5. For Go 1.11 vendoring, set environmental variable `GO111MODULE=on`
6. `go run main.go` in backend directory


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
3. Navigate to [http://localhost:8080](http://localhost:8080)

### Use:

Look at the backend logs to see the initialization of each dog location ID. 
The IDs listed can be entered into the client and used to receive dogs that are currently in those locations.
