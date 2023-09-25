package models

import (
	"database/sql"
	"fmt"

	"github.com/gastrader/website/rand"
)

const (
	MinBytesPerToken = 32
)

type Session struct {
	ID        int
	UserID    int
	// Token is only set when creating a new session. When looking up a session, this will be left empty. Only store hash of session token in DB, cannot reverse it into a raw token.
	Token string
	TokenHash string
}

type SessionService struct {
	DB *sql.DB
	BytesPerToken int
}

func (ss *SessionService) Create(userID int) (*Session, error){
	//TO DO:  Create session token
	bytesPerToken := ss.BytesPerToken
	if bytesPerToken < MinBytesPerToken {
		bytesPerToken = MinBytesPerToken
	}
	
	token, err := rand.String(ss.BytesPerToken)
	if err != nil{
		return nil, fmt.Errorf("create: %w", err)
	}
	session := Session{
		UserID: userID,
		Token: token,
		 //HASH THE SESSION TOKEN
	}
	return &session, nil
}

func (ss *SessionService) User(token string) (*User, error){
	//TO DO: Implement
	return nil, nil
}