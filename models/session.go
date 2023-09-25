package models

import (
	"crypto/sha256"
	"database/sql"
	"encoding/base64"
	"fmt"

	"github.com/gastrader/website/rand"
)

const (
	MinBytesPerToken = 32
)

type Session struct {
	ID     int
	UserID int
	// Token is only set when creating a new session. When looking up a session, this will be left empty. Only store hash of session token in DB, cannot reverse it into a raw token.
	Token     string
	TokenHash string
}

type SessionService struct {
	DB            *sql.DB
	BytesPerToken int
}

func (ss *SessionService) Create(userID int) (*Session, error) {
	bytesPerToken := ss.BytesPerToken
	if bytesPerToken < MinBytesPerToken {
		bytesPerToken = MinBytesPerToken
	}

	token, err := rand.String(bytesPerToken)
	fmt.Printf("The token is: %s\n", token)
	if err != nil {
		return nil, fmt.Errorf("create bytes: %w", err)
	}
	session := Session{
		UserID:    userID,
		Token:     token,
		TokenHash: ss.hash(token),
	}
	row := ss.DB.QueryRow(`
		UPDATE sessions
		SET token_hash = $2
		WHERE user_id = $1
		RETURNING id;`, session.UserID, session.TokenHash)
	err = row.Scan(&session.ID)
	if err == sql.ErrNoRows {
		row := ss.DB.QueryRow(`
			INSERT INTO sessions(user_id, token_hash) 
			VALUES($1, $2) 
			RETURNING id;`, session.UserID, session.TokenHash)
		err = row.Scan(&session.ID)
	}
	fmt.Printf("ID is: %v, hash is:  %s\n", session.UserID, session.TokenHash)
	if err != nil {
		return nil, fmt.Errorf("createe: %w", err)
	}
	return &session, nil
}

func (ss *SessionService) User(token string) (*User, error) {
	//TO DO: Implement
	return nil, nil
}

func (ss *SessionService) hash(token string) string {
	tokenHash := sha256.Sum256([]byte(token))
	return base64.URLEncoding.EncodeToString(tokenHash[:])
}
