package repository

import (
	"database/sql"
	"sync"

	"github.com/naufaldymahas/movie-grpc-api/entity"
)

var movieLogrepository MovieLogInterface
var doOnceMovieLogRepository sync.Once

func InitMovieLogRepository(db *sql.DB) MovieLogInterface {
	doOnceMovieLogRepository.Do(func() {
		movieLogrepository = &MovieLogRepository{
			DB: db,
		}
	})

	return movieLogrepository
}

type MovieLogRepository struct {
	DB *sql.DB
}

type MovieLogInterface interface {
	Create(data *entity.MovieLog) error
}

func (repo *MovieLogRepository) Create(data *entity.MovieLog) error {
	if _, err := repo.DB.Exec("INSERT INTO movie_logs(event_type, params) VALUES(?,?)", data.EventType, data.Params); err != nil {
		return err
	}
	return nil
}
