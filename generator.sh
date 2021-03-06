#!/bin/bash

# go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
# go get -u github.com/grpc-ecosystem/go-grpc-middleware
# go get -u go.uber.org/zap
# go get -u github.com/grpc-ecosystem/go-grpc-middleware/logging/zap

protoc greet/greetpb/greet.proto --go-grpc_out=./greet --go_out=./greet
protoc calculator/calculatorpb/calculator.proto --go-grpc_out=./calculator --go_out=./calculator
protoc blog/blogpb/blog.proto --go-grpc_out=./blog --go_out=./blog