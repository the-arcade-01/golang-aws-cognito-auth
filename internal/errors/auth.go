package errors

import (
	"errors"
	"fmt"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrInvalidInput       = errors.New("invalid input")
	ErrAccountExists      = errors.New("account already exists")
	ErrServiceUnavailable = errors.New("service unavailable")
	ErrPasswordReset      = errors.New("password reset required")
	ErrInvalidCode        = errors.New("invalid confirmation code")
	ErrExpiredCode        = errors.New("expired confirmation code")
)

type AuthError struct {
	StatusCode int
	Err        error
	Message    string
}

func (e *AuthError) Error() string {
	if e.Message != "" {
		return e.Message
	}
	return e.Err.Error()
}

func (e *AuthError) Unwrap() error {
	return e.Err
}

func NewInvalidCredentialsError(detail string) *AuthError {
	msg := "Invalid credentials"
	if detail != "" {
		msg = fmt.Sprintf("%s: %s", msg, detail)
	}
	return &AuthError{
		StatusCode: 401,
		Err:        ErrInvalidCredentials,
		Message:    msg,
	}
}

func NewInvalidInputError(detail string) *AuthError {
	msg := "Invalid input"
	if detail != "" {
		msg = fmt.Sprintf("%s: %s", msg, detail)
	}
	return &AuthError{
		StatusCode: 400,
		Err:        ErrInvalidInput,
		Message:    msg,
	}
}

func NewAccountExistsError() *AuthError {
	return &AuthError{
		StatusCode: 409,
		Err:        ErrAccountExists,
		Message:    "Account already exists",
	}
}

func NewServiceUnavailableError(detail string) *AuthError {
	msg := "Service unavailable"
	if detail != "" {
		msg = fmt.Sprintf("%s: %s", msg, detail)
	}
	return &AuthError{
		StatusCode: 503,
		Err:        ErrServiceUnavailable,
		Message:    msg,
	}
}

func NewPasswordResetError() *AuthError {
	return &AuthError{
		StatusCode: 401,
		Err:        ErrPasswordReset,
		Message:    "Password reset required",
	}
}

func NewInvalidCodeError(detail string) *AuthError {
	msg := "Invalid confirmation code"
	if detail != "" {
		msg = fmt.Sprintf("%s: %s", msg, detail)
	}
	return &AuthError{
		StatusCode: 400,
		Err:        ErrInvalidCode,
		Message:    msg,
	}
}

func NewExpiredCodeError() *AuthError {
	return &AuthError{
		StatusCode: 400,
		Err:        ErrExpiredCode,
		Message:    "Confirmation code has expired",
	}
}
