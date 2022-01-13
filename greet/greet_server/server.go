package main

import (
	"context"
	"fmt"
	"grpc-go-course/greet/greetpb"
	"log"
	"net"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_zap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/reflection"
)

var (
	zapLogger  *zap.Logger
	customFunc grpc_zap.CodeToLevel
)

func main() {
	fmt.Print("Hello world\n\n\n")

	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v\n", err)
	}

	opts := []grpc.ServerOption{}
	tls := true
	if tls {
		certFile := "ssl/private/server.crt"
		keyFile := "ssl/private/server.pem"
		creds, sslErr := credentials.NewServerTLSFromFile(certFile, keyFile)
		if sslErr != nil {
			log.Fatalf("Failed loading certificates: %v", sslErr)
			return
		}
		opts = append(opts, grpc.Creds(creds))
	}

	zapLogger, _ := zap.NewProduction()
	customFunc = grpc_zap.DefaultCodeToLevel
	alwaysLoggingDeciderServer := func(
		ctx context.Context,
		fullMethodName string,
		servingObject interface{},
	) bool {
		return true
	}

	// zapOpts := []grpc_zap.Option{
	// 	grpc_zap.WithLevels(customFunc),
	// }

	grpc_zap.ReplaceGrpcLoggerV2(zapLogger)

	opts = append(opts,
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			// grpc_ctxtags.UnaryServerInterceptor(
			// 	grpc_ctxtags.WithFieldExtractor(
			// 		grpc_ctxtags.CodeGenRequestFieldExtractor),
			// ),
			// grpc_zap.UnaryServerInterceptor(zapLogger, zapOpts...),
			grpc_zap.PayloadUnaryServerInterceptor(
				zapLogger,
				alwaysLoggingDeciderServer,
			),
		)))

	opts = append(opts,
		grpc.StreamInterceptor(grpc_middleware.ChainStreamServer(
			// grpc_ctxtags.StreamServerInterceptor(
			// 	grpc_ctxtags.WithFieldExtractor(
			// 		grpc_ctxtags.CodeGenRequestFieldExtractor),
			// ),
			// grpc_zap.StreamServerInterceptor(zapLogger, zapOpts...),
			grpc_zap.PayloadStreamServerInterceptor(
				zapLogger,
				alwaysLoggingDeciderServer,
			),
		)))

	s := grpc.NewServer(opts...)

	greetpb.RegisterGreetServiceServer(s, &server{})

	// Register reflection service on gRPC server.
	reflection.Register(s)

	if err := s.Serve(lis); err != nil {
		log.Fatalf("faild to serve %v\n", err)
	}

}
