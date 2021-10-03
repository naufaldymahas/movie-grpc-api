package server

import (
	"context"
	"log"
	"net"
	"os"
	"testing"

	"github.com/go-resty/resty/v2"
	"github.com/naufaldymahas/movie-grpc-api/pb"
	"github.com/naufaldymahas/movie-grpc-api/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc"
)

type RESTServerTestSuite struct {
	suite.Suite
	apiClient *resty.Client
}

func TestRESTServerTestSuite(t *testing.T) {
	suite.Run(t, new(RESTServerTestSuite))
}

func setupServer() *grpc.Server {
	ctx := context.Background()
	grpcPort := "8080"
	listener, _ := net.Listen("tcp", ":"+grpcPort)
	s := grpc.NewServer()
	svc := service.InitMovieService()
	pb.RegisterMovieServiceServer(s, svc)

	go func() {
		RESTServer(ctx, grpcPort, "8081")
	}()

	go func() {
		log.Println("Start GRPC Server")
		s.Serve(listener)
	}()

	return s
}

func (suite *RESTServerTestSuite) SetupTest() {
	os.Setenv("MYSQL_DB_USER", "root")
	os.Setenv("MYSQL_DB_PASSWORD", "password")
	os.Setenv("MYSQL_DB_NAME", "bibit-test")
	os.Setenv("OMDB_API", "http://www.omdbapi.com/")
	os.Setenv("OMDB_API_KEY", "faf7e5bb")
	suite.apiClient = resty.New()
}

func (suite *RESTServerTestSuite) TestRestServerSearchMovie() {
	s := setupServer()

	tests := []struct {
		statusCode int
		params     map[string]string
	}{
		{
			statusCode: 200,
			params: map[string]string{
				"pagination": "1",
				"searchword": "test",
			},
		},
		{
			statusCode: 400,
			params: map[string]string{
				"pagination": "",
				"searchword": "",
			},
		},
		{
			statusCode: 500,
			params:     map[string]string{},
		},
	}

	for _, tt := range tests {

		resp, err := suite.apiClient.R().
			SetQueryParams(tt.params).
			Get("http://localhost:8081/v1/movie")

		assert.NoError(suite.T(), err)
		assert.Equal(suite.T(), resp.StatusCode(), tt.statusCode)
	}

	s.GracefulStop()
}

func (suite *RESTServerTestSuite) TestRestServerSearchMovieByID() {
	s := setupServer()

	tests := []struct {
		statusCode int
		id         string
	}{
		{
			statusCode: 200,
			id:         "tt4244162",
		},
		{
			statusCode: 500,
		},
	}

	for _, tt := range tests {
		resp, err := suite.apiClient.R().
			Get("http://localhost:8081/v1/movie/" + tt.id)

		assert.NoError(suite.T(), err)
		assert.Equal(suite.T(), resp.StatusCode(), tt.statusCode)
	}

	s.GracefulStop()
}
