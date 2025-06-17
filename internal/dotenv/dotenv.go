package dotenv

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Env struct {
	PolkaKey  string
	DBUrl     string
	JWTSecret string
	Platform  string
}

func LoadEnv() (*Env, error) {
	godotenv.Load(".env")

	polkaKey, err := GetEnv("POLKA_KEY")
	if err != nil {
		return nil, fmt.Errorf("failed to load POLKA_KEY: %w", err)
	}

	dbURL, err := GetEnv("DB_URL")
	if err != nil {
		return nil, fmt.Errorf("failed to load DATABASE_URL: %w", err)
	}

	jwtSecret, err := GetEnv("JWT_SECRET")
	if err != nil {
		return nil, fmt.Errorf("failed to load JWT_TOKEN: %w", err)
	}

	platform, err := GetEnv("PLATFORM")
	if err != nil {
		return nil, fmt.Errorf("failed to load PLATFORM: %w", err)
	}

	return &Env{
		PolkaKey:  polkaKey,
		DBUrl:     dbURL,
		JWTSecret: jwtSecret,
		Platform:  platform,
	}, nil
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
