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

type UserinfoData struct {
	// can leave out fields - we are using a tolerant reader
	Audience      []string `json:"aud"`
	AuthTime      int64    `json:"auth_time"`
	Email         string   `json:"email"`
	EmailVerified bool     `json:"email_verified"`
	Name          string   `json:"name"` // username
	Groups        []string `json:"groups"`
	Issuer        string   `json:"iss"`
	IssuedAt      int64    `json:"iat"`
	RequestedAt   int64    `json:"rat"`
	Subject       string   `json:"sub"` //
}

type TokenIntrospectionData struct {
	Active    bool     `json:"active"`
	Scope     string   `json:"scope"`
	ClientId  string   `json:"client_id"`
	Sub       string   `json:"sub"`
	Exp       int64    `json:"exp"`
	Iat       int64    `json:"iat"`
	Nbf       int64    `json:"nbf"`
	Aud       []string `json:"aud"`
	Iss       string   `json:"iss"`
	TokenType string   `json:"token_type"`
	TokenUse  string   `json:"token_use"`

	// in case of error, you get these fields instead
	ErrorMessage string              `json:"message"`
	Errors       map[string][]string `json:"errors"`
}

type UserinfoResponseDto struct {
	UserinfoData

	// TODO temporarily retain old version
	Data UserinfoData `json:"data"`

	// in case of error, you get these fields instead (old version)
	ErrorCode        string `json:"error"`
	ErrorDescription string `json:"error_description"`
}

type IdentityProviderClient interface {
	TokenWithAuthenticationCodeAndPKCE(ctx context.Context, applicationConfigName string, authorizationCode string, pkceVerifier string) (*TokenResponseDto, int, error)

	UserInfo(ctx context.Context) (*UserinfoData, int, error)

	TokenIntrospection(ctx context.Context) (*TokenIntrospectionData, int, error)
}
