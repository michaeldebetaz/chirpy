package config

import (
	"database/sql"
	"fmt"

	"github.com/michaeldebetaz/chirpy/internal/database"
	"github.com/michaeldebetaz/chirpy/internal/dotenv"
	"github.com/michaeldebetaz/chirpy/internal/middlewares"
)

type Config struct {
	Env     *dotenv.Env
	Mw      *middlewares.Middleware
	Queries *database.Queries
}

func Init() (*Config, error) {
	env, err := dotenv.LoadEnv()
	if err != nil {
		err := fmt.Errorf("error loading environment variables: %v", err)
		return nil, err
	}

	db, err := sql.Open("postgres", env.DB_URL)
	if err != nil {
		err := fmt.Errorf("error connecting to the database: %v", err)
		return nil, err
	}

	return &Config{
		Env:     env,
		Mw:      middlewares.New(),
		Queries: database.New(db),
	}, nil
}
