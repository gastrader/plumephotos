package models

import (
	"crypto/sha256"
	"database/sql"
	"encoding/base64"
	"fmt"
	"strings"
	"time"

	"github.com/gastrader/website/rand"
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
	DB            *sql.DB
	BytesPerToken int
	Duration      time.Duration
}

func (pws *PasswordResetService) Create(email string) (*PasswordReset, error) {
	//Verify we have valid email address, get their ID
	email = strings.ToLower(email)
	var userID int
	row := pws.DB.QueryRow(`
		SELECT id FROM users WHERE email = $1;`, email)
	err := row.Scan(&userID)
	if err != nil {
		return nil, fmt.Errorf("create: %w", err)
	}
	//Build the password reset
	bytesPerToken := pws.BytesPerToken
	if bytesPerToken < MinBytesPerToken {
		bytesPerToken = MinBytesPerToken
	}
	token, err := rand.String(bytesPerToken)
	if err != nil {
		return nil, fmt.Errorf("create bytes: %w", err)
	}
	//Hash the token
	duration := pws.Duration
	if duration == 0 {
		duration = DefaultResetDuration
	}
	pwReset := PasswordReset{
		UserID:    userID,
		Token:     token,
		TokenHash: pws.hash(token),
		ExpiresAt: time.Now().Add(duration),
	}

	//Insert pwReset into DB
	row = pws.DB.QueryRow(`INSERT INTO password_resets (user_id, token_hash, expires_at) VALUES ($1, $2, $3) ON CONFLICT (user_id) DO UPDATE SET token_hash = $2, expires_at = $3 RETURNING id;`, pwReset.UserID, pwReset.TokenHash, pwReset.ExpiresAt)
	err = row.Scan(&pwReset.ID)
	if err != nil {
		return nil, fmt.Errorf("createe: %w", err)
	}

	return &pwReset, nil
}

func (pws *PasswordResetService) Consume(token string) (*User, error) {
	//verify pw reset is not expired
	tokenHash := pws.hash(token)
	fmt.Println("The token is: ", tokenHash)
	var user User
	var pwReset PasswordReset
	row := pws.DB.QueryRow(`
		SELECT password_resets.id, 
		password_resets.expires_at, 
		users.id, users.email, users.password_hash 
		FROM password_resets 
		JOIN users ON users.id = password_resets.user_id
		WHERE password_resets.token_hash = $1;`, tokenHash)
	err := row.Scan(&pwReset.ID, &pwReset.ExpiresAt,
		&user.ID, &user.Email, &user.PasswordHash)
	if err != nil {
		return nil, fmt.Errorf("consumee: %w", err)
	}
	if time.Now().After(pwReset.ExpiresAt) {
		return nil, fmt.Errorf("token expired: %v", token)
	}
	//delete pw reset.
	err = pws.delete(pwReset.ID)
	if err != nil {
		return nil, fmt.Errorf("consume: %w", err)
	}
	return &user, nil
}

func (pws *PasswordResetService) hash(token string) string {
	tokenHash := sha256.Sum256([]byte(token))
	return base64.URLEncoding.EncodeToString(tokenHash[:])
}

func (pws *PasswordResetService) delete(id int) error {
	_, err := pws.DB.Exec(`
	DELETE FROM password_resets WHERE id = $1;`, id)
	if err != nil {
		return fmt.Errorf("delete: %w", err)
	}
	return nil
}
