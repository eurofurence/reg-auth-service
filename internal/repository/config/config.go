package config

import (
	"crypto/rsa"
	"fmt"
	"time"
)

func UseEcsLogging() bool {
	return ecsLogging
}

func ServerAddr() string {
	c := configuration()
	return fmt.Sprintf("%s:%s", c.Server.Address, c.Server.Port)
}

func ServerReadTimeout() time.Duration {
	return time.Second * time.Duration(configuration().Server.ReadTimeout)
}

func ServerWriteTimeout() time.Duration {
	return time.Second * time.Duration(configuration().Server.WriteTimeout)
}

func ServerIdleTimeout() time.Duration {
	return time.Second * time.Duration(configuration().Server.IdleTimeout)
}

func IsCorsDisabled() bool {
	return configuration().Security.Cors.DisableCors
}

func CorsAllowOrigin() string {
	return configuration().Security.Cors.AllowOrigin
}

func SendInsecureCookies() bool {
	return configuration().Security.Cors.InsecureCookies
}

func SendNonHttpOnlyCookies() bool {
	return configuration().Security.Cors.DisableHttpOnlyCookies
}

func TokenEndpoint() string {
	return configuration().IdentityProvider.TokenEndpoint
}

func AuthorizationEndpoint() string {
	return configuration().IdentityProvider.AuthorizationEndpoint
}

func DropoffEndpointUrl() string {
	return configuration().Service.DropoffEndpointUrl
}

func TokenRequestTimeout() time.Duration {
	return configuration().IdentityProvider.TokenRequestTimeout
}

func AuthRequestTimeout() time.Duration {
	return configuration().IdentityProvider.AuthRequestTimeout
}

func GetApplicationConfig(applicationName string) (ApplicationConfig, error) {
	appConfig, found := configuration().ApplicationConfigs[applicationName]
	if found {
		return appConfig, nil
	} else {
		return ApplicationConfig{}, fmt.Errorf("no application configured for applicationName %s", applicationName)
	}
}

func LoggingSeverity() string {
	return configuration().Logging.Severity
}

func OidcIdTokenCookieName() string {
	return configuration().Security.Oidc.IdTokenCookieName
}

func OidcAccessTokenCookieName() string {
	return configuration().Security.Oidc.AccessTokenCookieName
}

func OidcKeySet() []*rsa.PublicKey {
	return parsedKeySet
}

func OidcAllowedAudience() string {
	return configuration().Security.Oidc.Audience
}

func OidcAllowedIssuer() string {
	return configuration().Security.Oidc.Issuer
}

func OidcUserInfoURL() string {
	return configuration().Security.Oidc.UserInfoURL
}

func OidcTokenIntrospectionURL() string {
	return configuration().Security.Oidc.TokenIntrospectionURL
}

func OidcUserInfoCacheRetentionTime() time.Duration {
	return time.Duration(configuration().Security.Oidc.UserInfoCacheSeconds) * time.Second
}

func OidcUserInfoCacheEnabled() bool {
	return configuration().Security.Oidc.UserInfoCacheSeconds > 0 &&
		configuration().Security.Oidc.UserInfoURL != "" &&
		configuration().Security.Oidc.AccessTokenCookieName != ""
}

func RelevantGroups() map[string][]string {
	return configuration().Security.Oidc.RelevantGroups
}
