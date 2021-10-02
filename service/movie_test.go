package service

import (
	"context"
	"log"
	"net"
	"os"
	"testing"

	"github.com/naufaldymahas/movie-grpc-api/pb"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
)

type MovieTestSuite struct {
	suite.Suite
	svc *MovieService
}

func dialer() func(context.Context, string) (net.Conn, error) {
	listener := bufconn.Listen(1024 * 1024)

	server := grpc.NewServer()

	pb.RegisterMovieServiceServer(server, InitMovieLogService())

	go func() {
		if err := server.Serve(listener); err != nil {
			log.Fatal(err)
		}
	}()

	return func(context.Context, string) (net.Conn, error) {
		return listener.Dial()
	}
}

func (suite *MovieTestSuite) SetupTest() {
	os.Setenv("MYSQL_DB_USER", "root")
	os.Setenv("MYSQL_DB_PASSWORD", "password")
	os.Setenv("MYSQL_DB_NAME", "bibit-test")
	os.Setenv("OMDB_API", "http://www.omdbapi.com/")
	os.Setenv("OMDB_API_KEY", "faf7e5bb")
	suite.svc = &MovieService{}
}

func (suite *MovieTestSuite) TestMovieCallRestPositive() {
	params := map[string]string{
		"s":    "Batman",
		"page": "1",
	}

	resp, err := suite.svc.callRest(params)
	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), 200, resp.StatusCode())
}

func (suite *MovieTestSuite) TestMovieSearchMoviePositive() {
	request := pb.FindAllRequest{
		Searchword: "batman",
		Pagination: 1,
	}

	result, err := suite.svc.SearchMovie(context.Background(), &request)
	assert.Nil(suite.T(), err)
	assert.NotNil(suite.T(), result)

	assert.Greater(suite.T(), result.TotalResult, int64(0))
	assert.NotEmpty(suite.T(), result.Results)
}

func (suite *MovieTestSuite) TestMovieSearchMovieNegative() {
	request := pb.FindAllRequest{}

	result, err := suite.svc.SearchMovie(context.Background(), &request)
	assert.NotNil(suite.T(), err)
	assert.Nil(suite.T(), result)
}

func (suite *MovieTestSuite) TestMovieSearchMovieByIDPositive() {
	request := pb.FindByIDRequest{
		ID: "tt0147746",
	}

	result, err := suite.svc.SearchMovieByID(context.Background(), &request)
	assert.Nil(suite.T(), err)
	assert.NotNil(suite.T(), result)

	assert.NotEqual(suite.T(), result.Title, "")
}

func (suite *MovieTestSuite) TestMovieSearchMovieByIDNegative_1() {
	request := pb.FindByIDRequest{
		ID: "tt0000000",
	}

	result, err := suite.svc.SearchMovieByID(context.Background(), &request)
	assert.NotNil(suite.T(), err)
	assert.Nil(suite.T(), result)
}

func (suite *MovieTestSuite) TestMovieSearchMovieByIDNegative_2() {
	request := pb.FindByIDRequest{
		ID: "1",
	}

	result, err := suite.svc.SearchMovieByID(context.Background(), &request)
	assert.NotNil(suite.T(), err)
	assert.Nil(suite.T(), result)
}

func (suite *MovieTestSuite) TestMovieSearchWithClient() {
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
			assert.Nil(suite.T(), err)
			assert.NotNil(suite.T(), result)

			assert.NotEmpty(suite.T(), result.Results)
		} else {
			assert.Nil(suite.T(), result)
			assert.NotNil(suite.T(), err)
		}

		if tt.totalMustGreaterThanZero {
			assert.Greater(suite.T(), result.TotalResult, 0)
		}
	}
}

func (suite *MovieTestSuite) TestMovieSearchByIDWithClient() {
	tests := []struct {
		req        *pb.FindByIDRequest
		isPositive bool
	}{
		{
			req: &pb.FindByIDRequest{
				ID: "tt4244162",
			},
			isPositive: true,
		},
		{
			req: &pb.FindByIDRequest{
				ID: "1",
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
			assert.Nil(suite.T(), err)
			assert.NotNil(suite.T(), result)
			assert.NotEqual(suite.T(), "", result.Title)
		} else {
			assert.Nil(suite.T(), result)
			assert.NotNil(suite.T(), err)
		}
	}
}

func TestMovieTestSuite(t *testing.T) {
	suite.Run(t, new(MovieTestSuite))
}
