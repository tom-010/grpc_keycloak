#!/bin/bash

protoc \
    --go_out=. \
    --go_opt=paths=source_relative \
    --go-grpc_out=. \
    --go-grpc_opt=paths=source_relative \
    proto/users.proto

protoc \
    --go_out=. \
    --go_opt=paths=source_relative \
    --go-grpc_out=. \
    --go-grpc_opt=paths=source_relative \
    proto/login.proto

# protoc \
#     --dart_out=grpc:fe/dart/lib/src/gen \
#     -Iproto proto/users.proto