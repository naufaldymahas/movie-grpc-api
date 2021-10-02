package repository

import (
	"fmt"
	"os"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/naufaldymahas/movie-grpc-api/config"
	"github.com/naufaldymahas/movie-grpc-api/entity"
	"github.com/stretchr/testify/assert"
)

func setDBEnv() {
	os.Setenv("MYSQL_DB_USER", "root")
	os.Setenv("MYSQL_DB_PASSWORD", "password")
	os.Setenv("MYSQL_DB_NAME", "bibit-test")
}

func TestMovieLogCreate(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.Nil(t, err)
	defer db.Close()

	repo := InitMovieLogRepository(db)
	data := entity.MovieLog{
		EventType: "test1",
		Params:    "test2",
	}

	mock.ExpectExec("INSERT INTO movie_logs").WithArgs(data.EventType, data.Params).WillReturnResult(sqlmock.NewResult(1, 1))
	err = repo.Create(&data)
	assert.Nil(t, err, fmt.Sprintf("Got error %v", err))

	err = mock.ExpectationsWereMet()
	assert.Nil(t, err, fmt.Sprintf("there were unfulfilled expectations: %s", err))
}

func TestMovieLogCreate_RealConnectionLoop5(t *testing.T) {
	setDBEnv()
	db := config.GetConn()
	assert.NotNil(t, db)

	repo := InitMovieLogRepository(db)
	data := []entity.MovieLog{
		{
			EventType: "loop 1",
			Params:    "loop 1",
		},
		{
			EventType: "loop 2",
			Params:    "loop 2",
		},
		{
			EventType: "loop 3",
			Params:    "loop 3",
		},
		{
			EventType: "loop 4",
			Params:    "loop 4",
		},
		{
			EventType: "loop 5",
			Params:    "loop 5",
		},
	}

	var err error
	for _, val := range data {
		err = repo.Create(&val)
		if err != nil {
			break
		}
	}
	assert.Nil(t, err, fmt.Sprintf("Got error %v", err))

}
