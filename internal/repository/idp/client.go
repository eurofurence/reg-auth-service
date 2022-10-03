package idp

import (
	"context"
	"fmt"
	aulogging "github.com/StephanHCB/go-autumn-logging"
	aurestbreaker "github.com/StephanHCB/go-autumn-restclient-circuitbreaker/implementation/breaker"
	aurestclientapi "github.com/StephanHCB/go-autumn-restclient/api"
	auresthttpclient "github.com/StephanHCB/go-autumn-restclient/implementation/httpclient"
	aurestlogging "github.com/StephanHCB/go-autumn-restclient/implementation/requestlogging"
	"github.com/eurofurence/reg-auth-service/internal/repository/config"
	"net/http"
	"net/url"
	"time"
)

type IdentityProviderClientImpl struct {
	client aurestclientapi.Client
}

// --- instance creation ---

func New() IdentityProviderClient {
	httpClient, err := auresthttpclient.New(0, nil, nil)
	if err != nil {
		aulogging.Logger.NoCtx().Fatal().WithErr(err).Printf("Failed to instantiate IDP client - BAILING OUT: %s", err.Error())
	}

	requestLoggingClient := aurestlogging.New(httpClient)

	circuitBreakerClient := aurestbreaker.New(requestLoggingClient,
		"identity-provider-breaker",
		10,
		2*time.Minute,
		30*time.Second,
		config.TokenRequestTimeout(),
	)

	return &IdentityProviderClientImpl{
		client: circuitBreakerClient,
	}
}

// --- implementation of repository interface ---

// can leave out fields to demo tolerant reader

func TokenRequestBody(appConfig config.ApplicationConfig, authorizationCode string, pkceVerifier string) url.Values {
	parameters := url.Values{}
	parameters.Set("grant_type", "authorization_code")
	parameters.Set("client_id", appConfig.ClientId)
	parameters.Set("client_secret", appConfig.ClientSecret)
	parameters.Set("redirect_uri", appConfig.DefaultDropoffUrl)
	parameters.Set("code", authorizationCode)
	parameters.Set("code_verifier", pkceVerifier)
	return parameters
}

func (i *IdentityProviderClientImpl) TokenWithAuthenticationCodeAndPKCE(ctx context.Context, applicationConfigName string, authorizationCode string, pkceVerifier string) (*TokenResponseDto, int, error) {
	appConfig, err := config.GetApplicationConfig(applicationConfigName)
	if err != nil {
		aulogging.Logger.Ctx(ctx).Warn().WithErr(err).Print(err.Error())
		return nil, http.StatusInternalServerError, err
	}

	requestBody := TokenRequestBody(appConfig, authorizationCode, pkceVerifier)
	tokenEndpoint := config.TokenEndpoint()
	bodyDto := TokenResponseDto{}
	response := aurestclientapi.ParsedResponse{
		Body: &bodyDto,
	}
	err = i.client.Perform(ctx, http.MethodPost, tokenEndpoint, requestBody, &response)
	if err != nil {
		aulogging.Logger.Ctx(ctx).Error().WithErr(err).Printf("error requesting token from identity provider: error from response is %s:%s, local error is %s", bodyDto.ErrorCode, bodyDto.ErrorDescription, err.Error())
		return nil, http.StatusBadGateway, err
	}
	if response.Status != http.StatusOK {
		err = fmt.Errorf("unexpected http status %d, was expecting %d", response.Status, http.StatusOK)
		aulogging.Logger.Ctx(ctx).Error().Printf("error requesting token from identity provider: error from response is %s:%s, local error is %s", bodyDto.ErrorCode, bodyDto.ErrorDescription, err.Error())
		return nil, response.Status, err
	}
	return &bodyDto, response.Status, nil
}
