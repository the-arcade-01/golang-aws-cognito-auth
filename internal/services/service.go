package services

import (
	"app/internal/models"
	"context"
)

type AuthServiceInterface interface {
	SignUp(ctx context.Context, user *models.User) (*models.DataResponse, *models.ErrorResponse)
	Login(ctx context.Context, user *models.UserLoginParams) (*models.DataResponse, *models.ErrorResponse)
	ConfirmAccount(ctx context.Context, user *models.UserConfirmationParams) (*models.DataResponse, *models.ErrorResponse)
	GetUser(ctx context.Context, token string) (*models.DataResponse, *models.ErrorResponse)
}
