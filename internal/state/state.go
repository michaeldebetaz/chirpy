package state

import (
	"database/sql"
	"fmt"

	"github.com/michaeldebetaz/chirpy/internal/database"
	"github.com/michaeldebetaz/chirpy/internal/dotenv"
	"github.com/michaeldebetaz/chirpy/internal/middlewares"
)

type State struct {
	Env     *dotenv.Env
	Mw      *middlewares.Middleware
	Db      *sql.DB
	Queries *database.Queries
}

func Init() (*State, error) {
	env, err := dotenv.LoadEnv()
	if err != nil {
		err := fmt.Errorf("Error loading environment variables: %v", err)
		return nil, err
	}

	db, err := sql.Open("postgres", env.DB_URL)
	if err != nil {
		err := fmt.Errorf("Error connecting to the database: %v", err)
		return nil, err
	}

	return &State{
		Env:     env,
		Mw:      middlewares.New(),
		Db:      db,
		Queries: database.New(db),
	}, nil
}
