package models

import (
	"database/sql"
	"fmt"

	_ "github.com/jackc/pgx/v5/stdlib"
)

// func DefaultPostgrsConfig() PostgresConfig {
// 	return PostgresConfig{
// 		Host:     "",
// 		Port:     "",
// 		User:     "",
// 		Password: "",
// 		Database: "",
// 		SSLMode:  "",
// 	}
// }

type PostgresConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Database string
	SSLMode  string
}

// Open opens a connection to the Postgres database.
// Callers of Open need to ensure that the connection
// is eventually closed by calling db.Close() method.
func Open(config PostgresConfig) (*sql.DB, error) {
	db, err := sql.Open("pgx", config.String())
	if err != nil {
		return nil, fmt.Errorf("open database: %w", err)
	}
	return db, nil
}

func (cfg PostgresConfig) String() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.Database, cfg.SSLMode)
}
