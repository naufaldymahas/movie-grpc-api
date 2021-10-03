package service

import (
	"database/sql"
	"encoding/json"
	"sync"

	"github.com/naufaldymahas/movie-grpc-api/entity"
	"github.com/naufaldymahas/movie-grpc-api/repository"
)

var movieLogService MovieLogInterface
var doOnceMovieLogService sync.Once

func InitMovieLogService(db *sql.DB) MovieLogInterface {
	doOnceMovieLogService.Do(func() {
		movieLogService = &MovieLogService{
			movieLogRepository: repository.InitMovieLogRepository(db),
		}
	})

	return movieLogService
}

type MovieLogService struct {
	movieLogRepository repository.MovieLogInterface
}

type MovieLogInterface interface {
	ToJSONString(s interface{}) string
	CreateMovieLog(eventType string, params string) error
}

func (svc *MovieLogService) ToJSONString(s interface{}) string {
	js, _ := json.Marshal(s)
	return string(js)
}

func (svc *MovieLogService) CreateMovieLog(eventType string, params string) error {
	data := entity.MovieLog{
		EventType: eventType,
		Params:    params,
	}

	return svc.movieLogRepository.Create(&data)
}
