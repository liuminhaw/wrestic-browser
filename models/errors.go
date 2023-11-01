package models

import "errors"

var (
	// A common pattern is to add the package as a prefix to the error for context.
	ErrUsernameTaken = errors.New("models: username is already in use")
	ErrNotFound      = errors.New("models: resource could not be found")
)
