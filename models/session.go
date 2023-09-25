package models

import "database/sql"

type Session struct {
	ID        int
	UserID    int
	// Token is only set when creating a new session. When looking up a session, this will be left empty. Only store hash of session token in DB, cannot reverse it into a raw token.
	Token string
	TokenHash string
}

type SessionService struct {
	DB *sql.DB
}

func (ss *SessionService) Create(userID int) (*Session, error){
	//TO DO:  Create session token
	return nil, nil
}

func (ss *SessionService) User(token string) (*User, error){
	//TO DO: Implement
	return nil, nil
}