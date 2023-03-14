package idp

import (
	"context"
	"fmt"
	aulogging "github.com/StephanHCB/go-autumn-logging"
	aurestbreaker "github.com/StephanHCB/go-autumn-restclient-circuitbreaker/implementation/breaker"
	aurestclientapi "github.com/StephanHCB/go-autumn-restclient/api"
	aurestcaching "github.com/StephanHCB/go-autumn-restclient/implementation/caching"
	auresthttpclient "github.com/StephanHCB/go-autumn-restclient/implementation/httpclient"
	aurestlogging "github.com/StephanHCB/go-autumn-restclient/implementation/requestlogging"
	"github.com/eurofurence/reg-auth-service/internal/repository/config"
	"github.com/eurofurence/reg-auth-service/internal/web/util/ctxvalues"
	"github.com/go-http-utils/headers"
	"net/http"
	"net/url"
	"time"
)

type IdentityProviderClientImpl struct {
	client aurestclientapi.Client
}

// --- instance creation ---

// useCacheCondition determines whether the cache should be used for a given request
//
// we cache only GET requests to the configured userinfo endpoint, and only for users who present a valid auth token
func useCacheCondition(ctx context.Context, method string, url string, requestBody interface{}) bool {
	return method == http.MethodGet && url == config.OidcUserInfoURL() && ctxvalues.AccessToken(ctx) != ""
}

// storeResponseCondition determines whether to store a response in the cache
//
// we only cache responses of successful requests to the userinfo endpoint
func storeResponseCondition(ctx context.Context, method string, url string, requestBody interface{}, response *aurestclientapi.ParsedResponse) bool {
	return response.Status == http.StatusOK
}

// cacheKeyFunction determines the key to cache the response under
//
// we cannot use the default cache key function, we must cache per auth token
func cacheKeyFunction(ctx context.Context, method string, requestUrl string, requestBody interface{}) string {
	return fmt.Sprintf("%s %s %s", ctxvalues.AccessToken(ctx), method, requestUrl)
}

// requestManipulator inserts Authorization when we are calling the userinfo endpoint
func requestManipulator(ctx context.Context, r *http.Request) {
	if r.Method == http.MethodGet {
		urlStr := r.URL.String()
		if urlStr != "" && (urlStr == config.OidcUserInfoURL() || urlStr == config.OidcTokenIntrospectionURL()) {
			r.Header.Set(headers.Authorization, "Bearer "+ctxvalues.AccessToken(ctx))
		}
	}
}

func New() IdentityProviderClient {
	httpClient, err := auresthttpclient.New(0, nil, requestManipulator)
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

	client := circuitBreakerClient

	if config.OidcUserInfoCacheEnabled() {
		cachingClient := aurestcaching.New(circuitBreakerClient,
			useCacheCondition,
			storeResponseCondition,
			cacheKeyFunction,
			config.OidcUserInfoCacheRetentionTime(),
			256,
		)
		client = cachingClient
	}

	return &IdentityProviderClientImpl{
		client: client,
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

func (i *IdentityProviderClientImpl) UserInfo(ctx context.Context) (*UserinfoData, int, error) {
	userinfoEndpoint := config.OidcUserInfoURL()
	bodyDto := UserinfoResponseDto{}
	response := aurestclientapi.ParsedResponse{
		Body: &bodyDto,
	}
	err := i.client.Perform(ctx, http.MethodGet, userinfoEndpoint, nil, &response)
	if err != nil {
		aulogging.Logger.Ctx(ctx).Error().WithErr(err).Printf("error requesting user info from identity provider: error from response is %s:%s, local error is %s", bodyDto.ErrorCode, bodyDto.ErrorDescription, err.Error())
		return nil, http.StatusBadGateway, err
	}
	if bodyDto.ErrorCode != "" || bodyDto.ErrorDescription != "" {
		aulogging.Logger.Ctx(ctx).Error().Printf("received an error response from identity provider: error from response is %s:%s", bodyDto.ErrorCode, bodyDto.ErrorDescription)
	}
	if response.Status != http.StatusOK && response.Status != http.StatusUnauthorized && response.Status != http.StatusForbidden {
		err = fmt.Errorf("unexpected http status %d, was expecting 200, 401, or 403", response.Status)
		aulogging.Logger.Ctx(ctx).Error().Printf("error requesting user info from identity provider: error from response is %s:%s, local error is %s", bodyDto.ErrorCode, bodyDto.ErrorDescription, err.Error())
		return nil, response.Status, err
	}
	if response.Status == http.StatusOK {
		if bodyDto.ErrorCode != "" || bodyDto.ErrorDescription != "" {
			err = fmt.Errorf("received an error response from identity provider: error from response is %s:%s", bodyDto.ErrorCode, bodyDto.ErrorDescription)
			return nil, response.Status, err
		}
	}

	if bodyDto.Data.Subject != "" {
		// got old response
		return &bodyDto.Data, response.Status, nil
	}

	return &bodyDto.UserinfoData, response.Status, nil
}

func (i *IdentityProviderClientImpl) TokenIntrospection(ctx context.Context) (*TokenIntrospectionData, int, error) {
	tokenIntrospectionEndpoint := config.OidcTokenIntrospectionURL()
	bodyDto := TokenIntrospectionData{}
	response := aurestclientapi.ParsedResponse{
		Body: &bodyDto,
	}
	err := i.client.Perform(ctx, http.MethodGet, tokenIntrospectionEndpoint, nil, &response)

	if err != nil {
		aulogging.Logger.Ctx(ctx).Error().WithErr(err).Printf("error requesting user info from identity provider: error from response is %s:%v, local error is %s", bodyDto.ErrorMessage, bodyDto.Errors, err.Error())
		return nil, http.StatusBadGateway, err
	}
	if bodyDto.ErrorMessage != "" || len(bodyDto.Errors) > 0 {
		aulogging.Logger.Ctx(ctx).Error().Printf("received an error response from identity provider: error from response is %s:%v", bodyDto.ErrorMessage, bodyDto.Errors)
	}
	if response.Status != http.StatusOK && response.Status != http.StatusUnauthorized && response.Status != http.StatusForbidden {
		err = fmt.Errorf("unexpected http status %d, was expecting 200, 401, or 403", response.Status)
		aulogging.Logger.Ctx(ctx).Error().Printf("error requesting user info from identity provider: error from response is %s:%v, local error is %s", bodyDto.ErrorMessage, bodyDto.Errors, err.Error())
		return nil, response.Status, err
	}
	if response.Status == http.StatusOK {
		if bodyDto.ErrorMessage != "" || len(bodyDto.Errors) > 0 {
			err = fmt.Errorf("received an error response from identity provider: error from response is %s:%v", bodyDto.ErrorMessage, bodyDto.Errors)
			return nil, response.Status, err
		}
	}

	return &bodyDto, response.Status, nil
}
