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
	Port:     "5432",
	User:     "postgres",
	Password: "admin123",
	Database: "website",
	
	}
}

type PostgresConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Database string
	
}

func (cfg PostgresConfig) String() string {
	 return fmt.Sprintf("postgres://%s:%s@%s:%s/%s", cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Database)
}