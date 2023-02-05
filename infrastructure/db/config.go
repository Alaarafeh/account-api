package db

import (
	"errors"
	"os"

	"github.com/joho/godotenv"
)

var errFailedToLoadEnv = errors.New("failed to load environment")

type Config struct {
	DBURL  string
	DBName string
}

func NewConfigFromEnvironment() (*Config, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, errFailedToLoadEnv
	}

	mongoDBURL := os.Getenv("MONGODB_URL")
	mongoDBName := os.Getenv("MONGODB_NAME")

	return &Config{DBURL: mongoDBURL, DBName: mongoDBName}, nil
}