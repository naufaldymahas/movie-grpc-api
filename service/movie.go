package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"sync"

	"github.com/go-resty/resty/v2"
	"github.com/jinzhu/copier"
	"github.com/naufaldymahas/movie-grpc-api/config"
	"github.com/naufaldymahas/movie-grpc-api/entity"
	"github.com/naufaldymahas/movie-grpc-api/pb"
	"github.com/naufaldymahas/movie-grpc-api/repository"
)

var movieService pb.MovieServiceServer
var doOnceMovieLogService sync.Once

func InitMovieLogService() pb.MovieServiceServer {
	doOnceMovieLogService.Do(func() {
		movieService = &MovieService{
			movieLogrepository: repository.InitMovieLogRepository(config.GetConn()),
		}
	})

	return movieService
}

type MovieService struct {
	movieLogrepository repository.MovieLogInterface
}

func (svc *MovieService) SearchMovie(ctx context.Context, req *pb.FindAllRequest) (*pb.FindAllResponse, error) {
	params := map[string]string{
		"s":    req.GetSearchword(),
		"page": fmt.Sprintf("%d", req.GetPagination()),
	}
	resp, err := svc.callRest(params)

	if resp.StatusCode() == 200 && err == nil {
		data := new(entity.MovieResponse)
		err = json.Unmarshal([]byte(resp.String()), data)

		if strings.EqualFold(data.Response, "false") {
			return nil, errors.New("failed to get movies")
		}

		if err == nil {
			result := new(pb.FindAllResponse)
			for _, val := range data.Search {
				movie := new(pb.MovieList)
				copier.Copy(movie, val)
				result.Results = append(result.Results, movie)
			}

			totalResult, _ := strconv.ParseInt(data.TotalResults, 10, 64)
			result.TotalResult = totalResult
			return result, nil
		}
	}
	return nil, err
}

func (svc *MovieService) SearchMovieByID(ctx context.Context, req *pb.FindByIDRequest) (*pb.Movie, error) {
	params := map[string]string{
		"i": req.GetID(),
	}

	resp, err := svc.callRest(params)

	if resp.StatusCode() == 200 && err == nil {
		data := new(entity.Movie)
		err = json.Unmarshal([]byte(resp.String()), data)

		if strings.EqualFold(data.Response, "false") {
			return nil, errors.New("failed to get movies")
		}

		if err == nil {
			result := new(pb.Movie)
			copier.Copy(result, data)
			return result, nil
		}
	}
	return nil, err
}

func (svc *MovieService) callRest(params map[string]string) (*resty.Response, error) {
	url := config.GetStringEnv("OMDB_API", "")
	client := resty.New()
	params["apikey"] = config.GetStringEnv("OMDB_API_KEY", "")

	return client.
		SetRetryCount(3).
		R().
		SetQueryParams(params).
		Get(url)
}
