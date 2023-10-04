package models

import "errors"

var (
	ErrEmailTaken = errors.New("models: Email address is already in use")
	ErrNotFound   = errors.New("models: Resource not found")
)
