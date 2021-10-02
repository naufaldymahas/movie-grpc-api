package config

import (
	"os"
	"testing"
)

func TestGetStringEnv(t *testing.T) {
	key := "TEST_KEY"
	expectationResult := "test"
	os.Setenv(key, expectationResult)

	value := GetStringEnv(key, "")
	if value != expectationResult {
		t.Errorf("Expectation result %s but get %s", expectationResult, value)
	}
}

func TestGetStringEnv_withDefaultKey(t *testing.T) {
	key := "TEST_KEY"
	expectationResult := "test"
	os.Setenv(key, expectationResult)

	value := GetStringEnv("forget_key", "test")
	if value != expectationResult {
		t.Errorf("Expectation result %s but get %s", expectationResult, value)
	}
}
