package acceptance

import (
	"context"
	"errors"
	"github.com/eurofurence/reg-auth-service/internal/web/util/ctxvalues"
	"net/http"
	"strings"

	"github.com/eurofurence/reg-auth-service/internal/repository/idp"
)

type mockIDPClient struct {
	recording []string
}

func (m *mockIDPClient) TokenWithAuthenticationCodeAndPKCE(ctx context.Context, applicationConfigName string, authorizationCode string, pkceVerifier string) (*idp.TokenResponseDto, int, error) {
	ret := &idp.TokenResponseDto{
		IdToken:     "dummy_mock_value",
		AccessToken: "access_mock_value",
	}
	return ret, http.StatusOK, nil
}

func (m *mockIDPClient) UserInfo(ctx context.Context) (*idp.UserinfoData, int, error) {
	ret := idp.UserinfoData{}

	token := ctxvalues.AccessToken(ctx)
	m.recording = append(m.recording, token)
	if token == "idp_is_down" {
		return &ret, http.StatusBadGateway, errors.New("simulated situation: idp unreachable")
	}
	if !strings.HasPrefix(token, "access_mock_value") {
		return &ret, http.StatusUnauthorized, nil
	}
	if token == "access_mock_value 101" {
		ret = idp.UserinfoData{
			Subject:       "101",
			Email:         "jsquirrel_github_9a6d@packetloss.de",
			EmailVerified: true,
			Name:          "me",
			Groups:        []string{"comedian", "fursuiter"},
		}
	} else if token == "access_mock_value 202" {
		ret = idp.UserinfoData{
			Subject:       "202",
			Email:         "jsquirrel_github_9a6d@packetloss.de",
			EmailVerified: true,
			Name:          "me",
			Groups:        []string{"comedian", "somethingelse"},
		}
	} else {
		ret = idp.UserinfoData{
			Subject:       "1234567890",
			Email:         "jsquirrel_github_9a6d@packetloss.de",
			EmailVerified: true,
			Name:          "me",
			Groups:        []string{"fursuiter", "staff", "admin"},
		}
	}
	return &ret, http.StatusOK, nil
}
