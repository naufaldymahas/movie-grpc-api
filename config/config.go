package config

import (
	"log"
	"os"
	"strconv"
	"sync"

	"github.com/joho/godotenv"
)

var doOnceConfig sync.Once

func InitConfig() {
	doOnceConfig.Do(func() {
		err := godotenv.Load()
		if err != nil {
			log.Fatal("Error load .env file")
		}
	})
}

func GetStringEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}

	return value
}

func GetIntEnv(key string, defaultValue int) int {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}

	valParse, _ := strconv.Atoi(value)
	return valParse
}

func GetBoolEnv(key string, defaultValue bool) bool {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}

	valParse, _ := strconv.ParseBool(value)
	return valParse
}
