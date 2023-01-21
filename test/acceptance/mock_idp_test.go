package acceptance

import (
	"context"
	"errors"
	"github.com/eurofurence/reg-auth-service/internal/web/util/ctxvalues"
	"net/http"

	"github.com/eurofurence/reg-auth-service/internal/repository/idp"
)

type mockIDPClient struct {
}

func (m *mockIDPClient) TokenWithAuthenticationCodeAndPKCE(ctx context.Context, applicationConfigName string, authorizationCode string, pkceVerifier string) (*idp.TokenResponseDto, int, error) {
	ret := &idp.TokenResponseDto{
		IdToken:     "dummy_mock_value",
		AccessToken: "access_mock_value",
	}
	return ret, http.StatusOK, nil
}

func (m *mockIDPClient) UserInfo(ctx context.Context) (*idp.UserinfoResponseDto, int, error) {
	token := ctxvalues.BearerAccessToken(ctx)
	if token == "Bearer idp_is_down" {
		return &idp.UserinfoResponseDto{}, http.StatusBadGateway, errors.New("simulated situation: idp unreachable")
	}
	if token != "Bearer access_mock_value" {
		return &idp.UserinfoResponseDto{}, http.StatusUnauthorized, nil
	}
	ret := &idp.UserinfoResponseDto{
		Subject: "1234567890",
		Global: idp.GlobalDto{
			Email:         "jsquirrel_github_9a6d@packetloss.de",
			EmailVerified: true,
			Name:          "me",
			Roles:         []string{"comedian", "fursuiter", "admin"},
		},
	}
	return ret, http.StatusOK, nil
}
