package idpclient

import (
	"context"
	"fmt"
	"github.com/eurofurence/reg-auth-service/internal/repository/config"
	"github.com/eurofurence/reg-auth-service/internal/repository/idp"
	"github.com/eurofurence/reg-auth-service/internal/repository/logging"
	"github.com/eurofurence/reg-auth-service/internal/repository/util/downstreamcall"
	"net/http"
	"time"
)

type TokenRequestDto struct {
	GrantType    string `json:"grant_type"`
	ClientId     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	// RedirectUri  string `json:"redirect_uri"`
	Code         string `json:"code"`
	CodeVerifier string `json:"code_verifier"`
}

type IdentityProviderClientImpl struct {
	netClient *http.Client
}

const HystrixCommandName = "idp_token"

// --- instance creation ---

func New() idp.IdentityProviderClient {
	// TODO lookup from configuration
	timeoutMs := 1000 // configuration.HttpRequestTimeoutMs()

	downstreamcall.ConfigureHystrixCommand(HystrixCommandName, int(timeoutMs))

	return &IdentityProviderClientImpl{
		netClient: &http.Client{
			// theoretically, this is no longer necessary with hystrix
			Timeout: time.Millisecond * time.Duration(timeoutMs) * 2,
		},
	}
}

// --- implementation of repository interface ---

// TODO - does the IDP have an error structure, and what is it like?

// intentionally leaving out fields to demo tolerant reader

type ErrorDto struct {
	Message   string `json:"message"`
}

func (i *IdentityProviderClientImpl) TokenWithAuthenticationCodeAndPKCE(ctx context.Context, applicationConfigName string, authorizationCode string, pkceVerifier string) (*idp.TokenResponseDto, error) {
	appConfig, err := config.GetApplicationConfig(applicationConfigName)
	if err != nil {
		return nil, err
	}
	requestDto := TokenRequestDto{
		GrantType: "authorization_code",
		ClientId: appConfig.ClientId,
		ClientSecret: appConfig.ClientSecret,
		// RedirectUri: config.RedirectUri(applicationConfig),
		Code: authorizationCode,
		CodeVerifier: pkceVerifier,
	}
	requestBody, err := downstreamcall.RenderJson(requestDto)
	if err != nil {
		return nil, err
	}

	tokenEndpoint := config.TokenEndpoint()

	responseBody, httpstatus, err := downstreamcall.HystrixPerformPOST(ctx, HystrixCommandName, i.netClient, tokenEndpoint, requestBody)

	if err != nil || httpstatus != http.StatusOK {
		if err == nil {
			err = fmt.Errorf("unexpected http status %d, was expecting %d", httpstatus, http.StatusOK)
		}

		errorResponseDto := &ErrorDto{}
		err2 := downstreamcall.ParseJson(responseBody, errorResponseDto)
		if err2 == nil {
			logging.Ctx(ctx).Error(fmt.Sprintf("error requesting token from identity provider: error from response is %s, local error is %s", errorResponseDto.Message, err.Error()))
		} else {
			logging.Ctx(ctx).Error(fmt.Sprintf("error requesting token from identity provider with no structured response available: local error is %s", err.Error()))
		}

		return nil, err
	}

	successResponseDto := &idp.TokenResponseDto{}
	err = downstreamcall.ParseJson(responseBody, successResponseDto)
	if err != nil {
		logging.Ctx(ctx).Error(fmt.Sprintf("error parsing token response from identity provider: error is %s", err.Error()))
		return nil, err
	}

	return successResponseDto, nil
}
