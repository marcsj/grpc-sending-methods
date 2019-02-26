#!/usr/bin/env bash
set -eu

KEY_DIR=../keys
BACKEND_DIR=../backend
WEB_DIR=../web-grpc/generated
GEN_TS_DIR=../web-grpc/node_modules/.bin/protoc-gen-ts


mkdir -p $KEY_DIR
mkdir -p $BACKEND_DIR
mkdir -p $WEB_DIR

# the following two lines are due to extensions generation in protoc-gen-ts
mkdir -p $WEB_DIR/google/api
touch $WEB_DIR/google/api/annotations_pb.js

GEN_NAME="dog"
mkdir -p $BACKEND_DIR/$GEN_NAME
protoc -I ./ \
  -I$GOPATH/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
  --swagger_out=logtostderr=true:$BACKEND_DIR/$GEN_NAME \
  --grpc-gateway_out=logtostderr=true:$BACKEND_DIR/$GEN_NAME \
  --plugin=protoc-gen-ts=$GEN_TS_DIR \
  --js_out=import_style=commonjs,binary:$WEB_DIR \
  --ts_out=service=true:$WEB_DIR \
  --go_out=plugins=grpc:$BACKEND_DIR/$GEN_NAME \
  $GEN_NAME.proto
