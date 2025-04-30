package services

import (
	"app/internal/db"
	appError "app/internal/errors"
	"app/internal/models"
	"context"
	"errors"
	"net/http"
)

type AuthService struct {
	store db.AuthStore
}

func NewAuthService(store db.AuthStore) *AuthService {
	return &AuthService{
		store: store,
	}
}

func (s *AuthService) SignUp(ctx context.Context, user *models.User) (*models.DataResponse, *models.ErrorResponse) {
	err := s.store.SignUp(ctx, user)
	if err != nil {
		var authErr *appError.AuthError
		if errors.As(err, &authErr) {
			return nil, models.NewErrorResponse(authErr.StatusCode, authErr.Error())
		}
		return nil, models.NewErrorResponse(http.StatusInternalServerError, "Failed to register user")
	}

	return models.NewDataResponse(http.StatusCreated, struct {
		Message string `json:"message"`
	}{
		Message: "User registered successfully.",
	}), nil
}

func (s *AuthService) Login(ctx context.Context, user *models.UserLoginParams) (*models.DataResponse, *models.ErrorResponse) {
	res, err := s.store.Login(ctx, user)
	if err != nil {
		var authErr *appError.AuthError
		if errors.As(err, &authErr) {
			return nil, models.NewErrorResponse(authErr.StatusCode, authErr.Error())
		}
		return nil, models.NewErrorResponse(http.StatusInternalServerError, err.Error())
	}

	return models.NewDataResponse(http.StatusOK, res), nil
}

func (s *AuthService) ConfirmAccount(ctx context.Context, user *models.UserConfirmationParams) (*models.DataResponse, *models.ErrorResponse) {
	err := s.store.ConfirmAccount(ctx, user)
	if err != nil {
		var authErr *appError.AuthError
		if errors.As(err, &authErr) {
			return nil, models.NewErrorResponse(authErr.StatusCode, authErr.Error())
		}
		return nil, models.NewErrorResponse(http.StatusInternalServerError, err.Error())
	}

	return models.NewDataResponse(http.StatusOK, struct {
		Message string `json:"message"`
	}{
		Message: "Account confirmed successfully.",
	}), nil
}

func (s *AuthService) GetUser(ctx context.Context, token string) (*models.DataResponse, *models.ErrorResponse) {
	res, err := s.store.GetUser(ctx, token)
	if err != nil {
		var authErr *appError.AuthError
		if errors.As(err, &authErr) {
			return nil, models.NewErrorResponse(authErr.StatusCode, authErr.Error())
		}
		return nil, models.NewErrorResponse(http.StatusInternalServerError, err.Error())
	}

	return models.NewDataResponse(http.StatusOK, res), nil
}
