#!/bin/bash

# PATH="${PATH}:${HOME}/go/bin" protoc greetpb/greet.proto --go-grpc_out=.

protoc --go_out=. --go_opt=paths=source_relative \
--go-grpc_out=. --go-grpc_opt=paths=source_relative \
greet.proto