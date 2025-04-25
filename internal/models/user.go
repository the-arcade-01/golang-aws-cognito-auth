package models

type User struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserLoginParams struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserConfirmationParams struct {
	Email string `json:"email"`
	Code  string `json:"code"`
}

type AuthLoginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"`
}

func NewAuthLoginResponse(accessToken, refreshToken string, expiresIn int) *AuthLoginResponse {
	return &AuthLoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    expiresIn,
	}
}
