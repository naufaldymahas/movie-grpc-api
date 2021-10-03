package service

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type MovieLogTestSuite struct {
	suite.Suite
}

func TestMovieLogTestSuite(t *testing.T) {
	suite.Run(t, new(MovieLogTestSuite))
}

func (suite *MovieLogTestSuite) TestMovieLogToJSONString() {
	sql, _, _ := sqlmock.New()
	svc := InitMovieLogService(sql)

	tests := []struct {
		testData     interface{}
		expectedData string
	}{
		{
			testData: map[string]string{
				"keyTest": "valueTest",
			},
			expectedData: `{"keyTest":"valueTest"}`,
		},
		{
			testData:     "",
			expectedData: `""`,
		},
	}

	for _, tt := range tests {
		result := svc.ToJSONString(tt.testData)
		assert.Equal(suite.T(), tt.expectedData, result)
	}
}

func (suite *MovieLogTestSuite) TestMovieLogCreateMovieLog() {
	sql, mock, _ := sqlmock.New()

	eventType := "test"
	params := "test"

	mock.ExpectExec("INSERT INTO movie_logs").WithArgs(eventType, params).WillReturnResult(sqlmock.NewResult(1, 1))

	svc := InitMovieLogService(sql)

	assert.NoError(suite.T(), svc.CreateMovieLog(eventType, params))
}
