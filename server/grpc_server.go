package server

import (
	"fmt"
	"net"

	"github.com/naufaldymahas/movie-grpc-api/pb"
	"github.com/naufaldymahas/movie-grpc-api/service"
	"google.golang.org/grpc"
)

func GRPCServer(port string) error {
	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		return err
	}

	s := grpc.NewServer()
	svc := service.InitMovieLogService()
	pb.RegisterMovieServiceServer(s, svc)

	fmt.Println("starting GRPC Server")

	return s.Serve(listener)
}
