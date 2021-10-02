package main

import (
	"log"

	"github.com/naufaldymahas/movie-grpc-api/config"
	"github.com/naufaldymahas/movie-grpc-api/server"
)

func main() {
	config.InitConfig()

	grpcPort := config.GetStringEnv("GRPC_PORT", "8080")

	grpcServer := server.GRPCServer(grpcPort)
	if grpcServer != nil {
		log.Fatal("GRPC Server failed to start")
	}
}
