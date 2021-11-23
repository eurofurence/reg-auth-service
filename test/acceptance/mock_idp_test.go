package acceptance

import (
	"context"

	"github.com/eurofurence/reg-auth-service/internal/repository/idp"
)

type mockIDPClient struct {
}


func (m *mockIDPClient) TokenWithAuthenticationCodeAndPKCE(ctx context.Context, applicationConfigName string, authorizationCode string, pkceVerifier string) (*idp.TokenResponseDto, error) {
	ret := &idp.TokenResponseDto{
		AccessToken:  "dummy_mock_value",
	}
	return ret, nil
}

