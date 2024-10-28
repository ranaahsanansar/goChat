package utils

// KeyCloakUser represents a user in Keycloak.
type KeyCloakUserDto struct {
	Username      string              `json:"username"`
	Enabled       bool                `json:"enabled"`
	Email         string              `json:"email"`
	FirstName     string              `json:"firstName"`
	LastName      string              `json:"lastName"`
	EmailVerified bool                `json:"emailVerified"`
	Credentials   []UserCredentialDto `json:"credentials"`
}

// UserCredential represents the user's credentials.
type UserCredentialDto struct {
	Type      string `json:"type"`
	Value     string `json:"value"`
	Temporary bool   `json:"temporary"`
}

type ApiError struct {
	ErrorMessage string `json:"errorMessage"`
}

type KeycloakLoginResponseDto struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type RefreshTokenResponse struct {
	AccessToken  string `json:"access_token"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
}
