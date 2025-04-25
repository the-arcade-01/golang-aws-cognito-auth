package db

import (
	"app/internal/config"
	appError "app/internal/errors"
	"app/internal/models"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"log/slog"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider/types"
)

type CognitoStore struct {
	client       *cognitoidentityprovider.Client
	userPoolId   string
	clientId     string
	clientSecret string
}

func NewCognitoStore(cfg *config.Config) *CognitoStore {
	store := &CognitoStore{
		userPoolId:   cfg.AwsCognitoUserPoolId,
		clientId:     cfg.AwsCognitoClientId,
		clientSecret: cfg.AwsCognitoClientSecret,
	}
	store.client = cognitoidentityprovider.NewFromConfig(cfg.AwsConfig)
	return store
}

func (s *CognitoStore) SignUp(ctx context.Context, user *models.User) error {
	input := &cognitoidentityprovider.SignUpInput{
		ClientId: aws.String(s.clientId),
		Password: aws.String(user.Password),
		Username: aws.String(user.Email),
		UserAttributes: []types.AttributeType{
			{Name: aws.String("name"), Value: aws.String(user.Name)},
			{Name: aws.String("email"), Value: aws.String(user.Email)},
		},
		SecretHash: aws.String(s.generateSecretHash(user.Email)),
	}

	_, err := s.client.SignUp(ctx, input)
	if err != nil {
		var usernameExistsErr *types.UsernameExistsException
		var invalidParamErr *types.InvalidParameterException
		var invalidPasswordErr *types.InvalidPasswordException

		if errors.As(err, &usernameExistsErr) {
			return appError.NewAccountExistsError()
		} else if errors.As(err, &invalidPasswordErr) {
			return appError.NewInvalidInputError(err.Error())
		} else if errors.As(err, &invalidParamErr) {
			return appError.NewInvalidInputError(err.Error())
		}

		slog.ErrorContext(ctx, "Failed to sign up user", "email", user.Email, "err", err)
		return appError.NewServiceUnavailableError("Unable to process registration")
	}

	return nil
}

func (s *CognitoStore) ConfirmAccount(ctx context.Context, user *models.UserConfirmationParams) error {
	_, err := s.client.ConfirmSignUp(ctx, &cognitoidentityprovider.ConfirmSignUpInput{
		ClientId:         aws.String(s.clientId),
		ConfirmationCode: aws.String(user.Code),
		Username:         aws.String(user.Email),
		SecretHash:       aws.String(s.generateSecretHash(user.Email)),
	})

	if err != nil {
		var codeMismatchErr *types.CodeMismatchException
		var expiredCodeErr *types.ExpiredCodeException
		var notFoundErr *types.UserNotFoundException

		switch {
		case errors.As(err, &codeMismatchErr):
			return appError.NewInvalidCodeError("")
		case errors.As(err, &expiredCodeErr):
			return appError.NewExpiredCodeError()
		case errors.As(err, &notFoundErr):
			return appError.NewInvalidInputError("User not found")
		default:
			slog.ErrorContext(ctx, "Failed to confirm account", "email", user.Email, "err", err)
			return appError.NewServiceUnavailableError("Unable to confirm account")
		}
	}

	return nil
}

func (s *CognitoStore) Login(ctx context.Context, user *models.UserLoginParams) (*models.AuthLoginResponse, error) {
	output, err := s.client.InitiateAuth(ctx, &cognitoidentityprovider.InitiateAuthInput{
		AuthFlow:       types.AuthFlowTypeUserPasswordAuth,
		ClientId:       aws.String(s.clientId),
		AuthParameters: map[string]string{"USERNAME": user.Email, "PASSWORD": user.Password, "SECRET_HASH": s.generateSecretHash(user.Email)},
	})

	if err != nil {
		var notAuthErr *types.NotAuthorizedException
		var userNotFoundErr *types.UserNotFoundException
		var userNotConfirmedErr *types.UserNotConfirmedException
		var passwordResetErr *types.PasswordResetRequiredException

		switch {
		case errors.As(err, &passwordResetErr):
			return nil, appError.NewPasswordResetError()
		case errors.As(err, &notAuthErr):
			return nil, appError.NewInvalidCredentialsError("")
		case errors.As(err, &userNotFoundErr):
			return nil, appError.NewInvalidCredentialsError("")
		case errors.As(err, &userNotConfirmedErr):
			return nil, appError.NewInvalidInputError("Account not confirmed")
		default:
			slog.ErrorContext(ctx, "Failed to sign in user", "email", user.Email, "err", err)
			return nil, appError.NewServiceUnavailableError("Authentication service unavailable")
		}
	}

	authResult := output.AuthenticationResult
	if authResult == nil {
		return nil, appError.NewServiceUnavailableError("Invalid authentication result")
	}

	return models.NewAuthLoginResponse(
		aws.ToString(authResult.AccessToken),
		aws.ToString(authResult.RefreshToken),
		int(authResult.ExpiresIn),
	), nil
}

func (s *CognitoStore) generateSecretHash(username string) string {
	h := hmac.New(sha256.New, []byte(s.clientSecret))
	h.Write([]byte(username + s.clientId))
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}
