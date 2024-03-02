package helper

import (
	"bytes"
)

const (
	Success             = "Success"
	ErrFailedReadBody   = "Failed to read body"
	ErrInternalError    = "Encountered an error"
	ErrFailParseRole    = "failed to get logged in role"
	ErrFailParsePayload = "failed to get logged in details"

	// Auth
	ErrFailHash           = "Failed to hash password"
	ErrUnauthorized       = "Action not allowed"
	ErrInvalidCredentials = "Invalid details"
	ErrExpiredAuthToken   = "Expired Token"

	// Token
	ErrGenerateToken       = "Cannot create confrimation code"
	ErrInvalidCode         = "Invalid code"
	ErrInvalidRefreshToken = "Invalid refresh token"
	ErrExpiredCode         = "Expired code"

	// User
	ErrExistingUsername   = "Existing username"
	ErrExistingEmail      = "Existing email"
	ErrNoExistingUsername = "No account with this username"
	ErrNoExistingEmail    = "No account with this email"

	// Mailer
	ErrCannotSendMail = "cannot send emails at the moment"
)

func FailGet(entity string) string {
	var buffer bytes.Buffer
	buffer.WriteString("failed to get ")
	buffer.Write([]byte(entity))
	return buffer.String()
}

func FailCreate(entity string) string {
	var buffer bytes.Buffer
	buffer.WriteString("failed to create ")
	buffer.Write([]byte(entity))
	return buffer.String()
}

func FailDelete(entity string) string {
	var buffer bytes.Buffer
	buffer.WriteString("failed to delete ")
	buffer.Write([]byte(entity))
	return buffer.String()
}

func NotFound(entity string) string {
	var buffer bytes.Buffer
	buffer.Write([]byte(entity))
	buffer.WriteString(" not found")
	return buffer.String()
}

func CodeSentToEmail(entity string) string {
	var buffer bytes.Buffer
	buffer.WriteString("Code sent to email ")
	buffer.Write([]byte(entity))
	return buffer.String()
}
