package idp

import "context"

// TODO - is this the correct response data structure of the IDP?

// can leave out fields - we are using a tolerant reader

type TokenResponseDto struct {
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	AccessToken  string `json:"access_token"`
	Scope        string `json:"scope"`
	RefreshToken string `json:"refresh_token"`
}

type IdentityProviderClient interface {
	TokenWithAuthenticationCodeAndPKCE(ctx context.Context, applicationConfigName string, authorizationCode string, pkceVerifier string) (*TokenResponseDto, error)
}
