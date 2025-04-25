package db

import (
	"app/internal/models"
	"context"
)

type AuthStore interface {
	SignUp(ctx context.Context, user *models.User) error
	ConfirmAccount(ctx context.Context, user *models.UserConfirmationParams) error
	Login(ctx context.Context, user *models.UserLoginParams) (*models.AuthLoginResponse, error)
}
