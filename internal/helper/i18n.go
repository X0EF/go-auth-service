package helper

import (
	"fmt"
)

const (
	Success             = "Success"
	ErrFailedReadBody   = "Failed to read body"
	ErrInternalError    = "Encountered an error"
	ErrFailParsePayload = "failed to get logged in details"

	// Auth
	ErrFailHash     = "Failed to hash password"
	ErrUnauthorized = "Action not allowed"
	ErrInvalidToken = "Invalid Token"
	ErrExpiredToken = "Expired Token"

	// User
	ErrExistingUsername   = "Existing username"
	ErrExistingEmail      = "Existing email"
	ErrNoExistingUsername = "No account with this username"
	ErrNoExistingEmail    = "No account with this email"
	ErrFailCreateUser     = "Failed to create user"
	ErrFailDelUser        = "failed to delete user"
	ErrFailGetUsers       = "Failed to retrieve users"
)

func NotFound(entity string) string {
	return fmt.Sprintf("%s not found", entity)
}
