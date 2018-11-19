#!/usr/bin/env bash
set -eu

mkdir -p ../keys
mkdir -p ../backend/dog
mkdir -p ../web/generated

protoc -I ./ \
  -I$GOPATH/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
  --swagger_out=logtostderr=true:../backend/dog \
  --grpc-gateway_out=logtostderr=true:../backend/dog \
  --plugin=protoc-gen-ts=../web/node_modules/.bin/protoc-gen-ts \
  --js_out=import_style=commonjs,binary:../web/generated \
  --ts_out=service=true:../web/generated \
  --go_out=plugins=grpc:../backend/dog \
  dog.proto
