package acceptance

import (
	"context"
	"net/http"

	"github.com/eurofurence/reg-auth-service/internal/repository/idp"
)

type mockIDPClient struct {
}


func (m *mockIDPClient) TokenWithAuthenticationCodeAndPKCE(ctx context.Context, applicationConfigName string, authorizationCode string, pkceVerifier string) (*idp.TokenResponseDto, int, error) {
	ret := &idp.TokenResponseDto{
		IdToken:  "dummy_mock_value",
	}
	return ret, http.StatusOK, nil
}
