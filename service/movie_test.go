package service

import (
	"context"
	"os"
	"testing"

	"github.com/naufaldymahas/movie-grpc-api/pb"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type MovieTestSuite struct {
	suite.Suite
	svc *MovieService
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
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), 200, resp.StatusCode())
}

func (suite *MovieTestSuite) TestMovieSearchMoviePositive() {
	request := pb.FindAllRequest{
		Searchword: "batman",
		Pagination: 1,
	}

	result, err := suite.svc.SearchMovie(context.Background(), &request)
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), result)

	assert.Greater(suite.T(), result.TotalResult, int64(0))
	assert.NotEmpty(suite.T(), result.Results)
}

func (suite *MovieTestSuite) TestMovieSearchMovieNegative() {
	request := pb.FindAllRequest{}

	result, err := suite.svc.SearchMovie(context.Background(), &request)
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), result)
}

func (suite *MovieTestSuite) TestMovieSearchMovieByIDPositive() {
	request := pb.FindByIDRequest{
		Id: "tt0147746",
	}

	result, err := suite.svc.SearchMovieByID(context.Background(), &request)
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), result)

	assert.NotEqual(suite.T(), result.Title, "")
}

func (suite *MovieTestSuite) TestMovieSearchMovieByIDNegative_1() {
	request := pb.FindByIDRequest{
		Id: "tt0000000",
	}

	result, err := suite.svc.SearchMovieByID(context.Background(), &request)
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), result)
}

func (suite *MovieTestSuite) TestMovieSearchMovieByIDNegative_2() {
	request := pb.FindByIDRequest{
		Id: "1",
	}

	result, err := suite.svc.SearchMovieByID(context.Background(), &request)
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), result)
}

func TestMovieTestSuite(t *testing.T) {
	suite.Run(t, new(MovieTestSuite))
}
