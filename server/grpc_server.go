package server

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"

	"github.com/naufaldymahas/movie-grpc-api/pb"
	"github.com/naufaldymahas/movie-grpc-api/service"
	"google.golang.org/grpc"
)

func GRPCServer(ctx context.Context, port string) error {
	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		return err
	}

	s := grpc.NewServer()
	svc := service.InitMovieService()
	pb.RegisterMovieServiceServer(s, svc)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for range c {
			log.Println("Shutting down GRPC Server")

			s.GracefulStop()

			<-ctx.Done()
		}
	}()

	log.Println("starting GRPC Server")
	return s.Serve(listener)
}
