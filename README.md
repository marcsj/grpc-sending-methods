# Dog Tracker
An example of grpc-web using a streaming dog tracker. This is a simple example of streaming grpc to grpc-web, and does not use an actual datastore apart from the list of dogs created in memory.

## Building:
1. Run `generate.sh` in the proto folder.
2. For Go 1.11 vendoring, set envionmental variable `GO111MODULE=on`.
3. Run `go mod vendor` to update vendor folder.
4. `go build` or `go run`.