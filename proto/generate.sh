#!/usr/bin/env bash

protoc -I ./ \
       --go_out=plugins=grpc:../backend/dog \
       dog.proto

# TODO: add js generation when we add clients:
#--js_out=import_style=commonjs,binary:../client/js/_proto \