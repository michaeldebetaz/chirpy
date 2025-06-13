package dotenv

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Env struct {
	DB_URL   string
	PLATFORM string
}

func LoadEnv() (*Env, error) {
	godotenv.Load(".env")

	dbURL, err := GetEnv("DB_URL")
	if err != nil {
		return nil, fmt.Errorf("failed to load DATABASE_URL: %w", err)
	}

	platform, err := GetEnv("PLATFORM")
	if err != nil {
		return nil, fmt.Errorf("failed to load PLATFORM: %w", err)
	}

	return &Env{DB_URL: dbURL, PLATFORM: platform}, nil
}

func GetEnv(key string) (string, error) {
	value, ok := os.LookupEnv(key)
	if !ok {
		err := fmt.Errorf("environment variable %s not set", key)
		return "", err
	}

	if value == "" {
		err := fmt.Errorf("environment variable %s is empty", key)
		return "", err
	}

	return value, nil
}
