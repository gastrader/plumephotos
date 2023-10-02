package models

import (
	"database/sql"
	"time"
)

const (
	DefaultResetDuration = 1 * time.Hour
)

type PasswordReset struct {
	ID     int
	UserID int
	//Token only set when passwordreset is being created.
	Token     string
	TokenHash string
	ExpiresAt time.Time
}

type PasswordResetService struct {
	DB *sql.DB
	BytesPerToken int
	Duration time.Duration
}

func (pws *PasswordResetService) Create(email string) (*PasswordReset, error){
	return nil, nil
}

func (pws *PasswordResetService) Consume(token string) (*User, error){
	return nil, nil 
}