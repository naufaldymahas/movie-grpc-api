package main

import (
	"context"

	"github.com/naufaldymahas/movie-grpc-api/config"
	"github.com/naufaldymahas/movie-grpc-api/server"
)

func main() {
	config.InitConfig()

	ctx := context.Background()

	grpcPort := config.GetStringEnv("GRPC_PORT", "8080")
	restPort := config.GetStringEnv("REST_PORT", "8081")

	go func() {
		server.RESTServer(ctx, grpcPort, restPort)
	}()

	server.GRPCServer(ctx, grpcPort)
}
