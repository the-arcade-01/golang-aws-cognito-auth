package db

import (
	"app/internal/models"
	"context"

	"github.com/golang-jwt/jwt/v5"
)

type AuthStore interface {
	ValidateToken(tokenString string) (*jwt.Token, error)
	GetClaims(token *jwt.Token) (map[string]interface{}, error)
	SignUp(ctx context.Context, user *models.User) error
	ConfirmAccount(ctx context.Context, user *models.UserConfirmationParams) error
	Login(ctx context.Context, user *models.UserLoginParams) (*models.AuthLoginResponse, error)
	GetUser(ctx context.Context, token string) (*models.UserInfoResponse, error)
}
