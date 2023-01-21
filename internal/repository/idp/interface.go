package idp

import "context"

type TokenResponseDto struct {
	// can leave out fields - we are using a tolerant reader
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
	IdToken     string `json:"id_token"`
	Scope       string `json:"scope"`
	TokenType   string `json:"token_type"`

	// in case of error, you get these fields instead
	ErrorCode        string `json:"error"`
	ErrorDescription string `json:"error_description"`
}

type UserinfoResponseDto struct {
	// can leave out fields - we are using a tolerant reader
	Audience    []string  `json:"aud"`
	AuthTime    int64     `json:"auth_time"`
	Global      GlobalDto `json:"global"`
	Issuer      string    `json:"iss"`
	IssuedAt    int64     `json:"iat"`
	RequestedAt int64     `json:"rat"`
	Subject     string    `json:"sub"` //

	// in case of error, you get these fields instead
	ErrorCode        string `json:"error"`
	ErrorDescription string `json:"error_description"`
}

type GlobalDto struct {
	Email         string   `json:"email"`
	EmailVerified bool     `json:"email_verified"`
	Name          string   `json:"name"` // username
	Roles         []string `json:"roles"`
}

type IdentityProviderClient interface {
	TokenWithAuthenticationCodeAndPKCE(ctx context.Context, applicationConfigName string, authorizationCode string, pkceVerifier string) (*TokenResponseDto, int, error)

	UserInfo(ctx context.Context) (*UserinfoResponseDto, int, error)
}
