# gRPC Sending Methods

Examples of streaming and otherwise sending gRPC from a go backend to clients using a dog tracker. 
Data is sent from the backend via these methods:
- **[gRPC](https://github.com/grpc/grpc-go)**
- **[gRPC-gateway](https://github.com/grpc-ecosystem/grpc-gateway)**(JSON over http)
- **[grpc-websocket-proxy](https://github.com/tmc/grpc-websocket-proxy)**(JSON over websockets), 
- **[grpc-web](https://github.com/improbable-eng/grpc-web)**

These sending methods enable you to generate interfaces and models with gRPC once,
 and to be called in multiple different ways. The most important way this helps 
 you leverage code is by streaming to a web client.
 
## Server:
`go run backend/cmd/server/main.go`


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

### OpenAPI
Per the example, our JSON documentation generates to `:8082/docs/{proto_name}`  
 [http://localhost:8082/docs/dog](http://localhost:8082/docs/dog)
### Use

Look at the backend logs to see the initialization of each dog location ID. 
The IDs listed can be entered into the client and used to receive dogs that are currently in those locations.
