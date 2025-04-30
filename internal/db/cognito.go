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
	"fmt"
	"log/slog"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider/types"
	"github.com/golang-jwt/jwt/v5"
	"github.com/lestrrat-go/jwx/v2/jwk"
)

type CognitoStore struct {
	client       *cognitoidentityprovider.Client
	userPoolId   string
	clientId     string
	clientSecret string
	tokenURL     string
	jwtIssuerURL string
	jwkSet       jwk.Set
}

func NewCognitoStore(cfg *config.Config) (*CognitoStore, error) {
	keySet, err := jwk.Fetch(context.Background(), cfg.AwsTokenURL)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch JWK set: %w", err)
	}
	return &CognitoStore{
		userPoolId:   cfg.AwsCognitoUserPoolId,
		clientId:     cfg.AwsCognitoClientId,
		clientSecret: cfg.AwsCognitoClientSecret,
		tokenURL:     cfg.AwsTokenURL,
		jwtIssuerURL: cfg.AwsJWTIssuerURL,
		client:       cognitoidentityprovider.NewFromConfig(cfg.AwsConfig),
		jwkSet:       keySet,
	}, nil
}

func (s *CognitoStore) ValidateToken(tokenString string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		kid, ok := token.Header["kid"].(string)
		if !ok {
			return nil, errors.New("key ID not found in token")
		}

		key, found := s.jwkSet.LookupKeyID(kid)
		if !found {
			return nil, errors.New("key not found in JWKS")
		}

		var rawKey interface{}
		if err := key.Raw(&rawKey); err != nil {
			return nil, fmt.Errorf("failed to get raw key: %w", err)
		}

		return rawKey, nil
	})
	if err != nil {
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token claims")
	}

	issuer, err := claims.GetIssuer()
	if err != nil {
		return nil, errors.New("token has invalid issuer")
	}

	if strings.Compare(issuer, s.jwtIssuerURL) != 0 {
		return nil, errors.New("token was not issued by the specified Cognito user pool")
	}

	return token, nil
}

func (s *CognitoStore) GetClaims(token *jwt.Token) (map[string]interface{}, error) {
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("invalid token claims")
	}
	return claims, nil
}

func (s *CognitoStore) GetUser(ctx context.Context, token string) (*models.UserInfoResponse, error) {
	output, err := s.client.GetUser(ctx, &cognitoidentityprovider.GetUserInput{
		AccessToken: aws.String(token),
	})

	if err != nil {
		var forbiddenErr *types.ForbiddenException
		var invalidParamErr *types.InvalidParameterException
		var notAuthErr *types.NotAuthorizedException

		if errors.As(err, &forbiddenErr) || errors.As(err, &invalidParamErr) || errors.As(err, &notAuthErr) {
			return nil, appError.NewInvalidInputError(err.Error())
		}

		slog.ErrorContext(ctx, "Failed to get user info", "err", err)
		return nil, appError.NewServiceUnavailableError("Unable to fetch user info")
	}

	attributes, username := output.UserAttributes, output.Username
	if attributes == nil || username == nil {
		return nil, appError.NewServiceUnavailableError("Invalid user info result")
	}

	attributesMap := make(map[string]string)
	for _, attribute := range attributes {
		attributesMap[aws.ToString(attribute.Name)] = aws.ToString(attribute.Value)
	}

	return &models.UserInfoResponse{
		Attributes: attributesMap,
		Username:   aws.ToString(username),
	}, nil
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
