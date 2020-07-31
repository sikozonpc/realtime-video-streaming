package postgresdb

import (
	"database/sql"
	"fmt"
	"goproject/env"
	"log"

	_ "github.com/lib/pq"
)

// New creates a new connection to the postgres database
func New(envVars env.Variables) (*sql.DB, error) {
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s sslmode=disable",
		envVars.DbHost,
		envVars.DbPort,
		envVars.DbUser,
		envVars.DbPassword,
		envVars.DbName,
	)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return nil, err
	}

	// Connection to database
	err = db.Ping()
	if err != nil {
		return nil, err
	}

	log.Printf("Connected to %s database", envVars.DbName)

	return db, nil
}
