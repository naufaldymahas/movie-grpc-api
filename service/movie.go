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
)

var movieService pb.MovieServiceServer
var doOnceMovieService sync.Once

func InitMovieService() pb.MovieServiceServer {
	doOnceMovieService.Do(func() {
		movieService = &MovieService{}
		InitMovieLogService(config.GetConn())
	})

	return movieService
}

type MovieService struct {
}

func (svc *MovieService) SearchMovie(ctx context.Context, req *pb.FindAllRequest) (*pb.FindAllResponse, error) {
	params := map[string]string{
		"s":    req.GetSearchword(),
		"page": fmt.Sprintf("%d", req.GetPagination()),
	}
	resp, err := svc.callRest(params)

	rawParam := movieLogService.ToJSONString(req)
	go movieLogService.CreateMovieLog("Search Movie", rawParam)
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
		"i": req.GetId(),
	}

	resp, err := svc.callRest(params)

	rawParam := movieLogService.ToJSONString(req)
	go movieLogService.CreateMovieLog("Search Movie By ID", rawParam)
	if resp.StatusCode() == 200 && err == nil {
		data := new(entity.Movie)
		err = json.Unmarshal([]byte(resp.String()), data)

		if strings.EqualFold(data.Response, "false") {
			return nil, errors.New("failed to get detail movie")
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
