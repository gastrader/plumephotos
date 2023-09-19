package models

import (
	"database/sql"
	"fmt"
	_ "github.com/jackc/pgx/v4/stdlib"
)

//WILL OPEN SQL CONNECTION WITH POSTGRES DB, ENSURE CONNECTION IS CLOSED db.Close() method
func Open(cfg PostgresConfig) (*sql.DB, error){
	db, err := sql.Open("pgx", cfg.String())
	if err !=nil {
		return nil, fmt.Errorf("open: %w", err)
	}
	return db, nil
}

func DefaultPostgresConfig() PostgresConfig{
	return  PostgresConfig{
	Host:    "localhost",
	Port:     "4321",
	User:     "postgres",
	Password: "admin123",
	Database: "website",
	SSLMode:  "disable",
	}
}

type PostgresConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Database string
	SSLMode  string
}

func (cfg PostgresConfig) String() string {
    return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s", cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.Database, cfg.SSLMode)
}