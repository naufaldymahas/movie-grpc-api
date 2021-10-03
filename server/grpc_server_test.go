package server

import (
	"context"
	"log"
	"net"
	"os"
	"testing"

	"github.com/naufaldymahas/movie-grpc-api/pb"
	"github.com/naufaldymahas/movie-grpc-api/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
)

type GRPCServerTestSuite struct {
	suite.Suite
}

func dialer() func(context.Context, string) (net.Conn, error) {
	listener := bufconn.Listen(1024 * 1024)

	server := grpc.NewServer()

	pb.RegisterMovieServiceServer(server, service.InitMovieService())

	go func() {
		if err := server.Serve(listener); err != nil {
			log.Fatal(err)
		}
	}()

	return func(context.Context, string) (net.Conn, error) {
		return listener.Dial()
	}
}

func TestGRPCServerTestSuite(t *testing.T) {
	suite.Run(t, new(GRPCServerTestSuite))
}

func (suite *GRPCServerTestSuite) SetupTest() {
	os.Setenv("MYSQL_DB_USER", "root")
	os.Setenv("MYSQL_DB_PASSWORD", "password")
	os.Setenv("MYSQL_DB_NAME", "bibit-test")
	os.Setenv("OMDB_API", "http://www.omdbapi.com/")
	os.Setenv("OMDB_API_KEY", "faf7e5bb")
}

func (suite *GRPCServerTestSuite) TestGRPCServerMovieSearchByID() {
	tests := []struct {
		req        *pb.FindByIDRequest
		isPositive bool
	}{
		{
			req: &pb.FindByIDRequest{
				Id: "tt4244162",
			},
			isPositive: true,
		},
		{
			req: &pb.FindByIDRequest{
				Id: "1",
			},
			isPositive: false,
		},
	}

	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, "", grpc.WithInsecure(), grpc.WithContextDialer(dialer()))
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	client := pb.NewMovieServiceClient(conn)

	for _, tt := range tests {
		result, err := client.SearchMovieByID(ctx, tt.req)

		if tt.isPositive {
			assert.NoError(suite.T(), err)
			assert.NotNil(suite.T(), result)
			assert.NotEqual(suite.T(), "", result.Title)
		} else {
			assert.Nil(suite.T(), result)
			assert.Error(suite.T(), err)
		}
	}
}

func (suite *GRPCServerTestSuite) TestGRPCServerMovieSearch() {
	tests := []struct {
		req                      *pb.FindAllRequest
		isPositive               bool
		totalMustGreaterThanZero bool
	}{
		{
			req: &pb.FindAllRequest{
				Searchword: "batman",
				Pagination: 1,
			},
			isPositive:               true,
			totalMustGreaterThanZero: true,
		},
		{
			req:                      &pb.FindAllRequest{},
			isPositive:               false,
			totalMustGreaterThanZero: false,
		},
	}

	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, "", grpc.WithInsecure(), grpc.WithContextDialer(dialer()))
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	client := pb.NewMovieServiceClient(conn)

	for _, tt := range tests {
		result, err := client.SearchMovie(ctx, tt.req)

		if tt.isPositive {
			assert.NoError(suite.T(), err)
			assert.NotNil(suite.T(), result)

			assert.NotEmpty(suite.T(), result.Results)
		} else {
			assert.Nil(suite.T(), result)
			assert.Error(suite.T(), err)
		}

		if tt.totalMustGreaterThanZero {
			assert.Greater(suite.T(), result.TotalResult, int64(0))
		}
	}
}
