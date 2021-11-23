package config

import (
	"testing"
	"time"

	"github.com/eurofurence/reg-auth-service/docs"
	"github.com/stretchr/testify/require"
)

func tstValidatePort(t *testing.T, value string, errMessage string) {
	errs := validationErrors{}
	config := serverConfig{Port: value}
	validateServerConfiguration(errs, config)
	require.Equal(t, 1, len(errs))
	require.Equal(t, []string{errMessage}, errs["server.port"])
}

func TestValidateServerConfiguration_empty(t *testing.T) {
	docs.Description("validation should catch an empty port configuration")
	tstValidatePort(t, "", "value '' cannot be empty")
}

func TestValidateServerConfiguration_numeric(t *testing.T) {
	docs.Description("validation should catch a non-numeric port configuration")
	tstValidatePort(t, "katze", "value 'katze' is not a valid port number")
}

func TestValidateServerConfiguration_tooHigh(t *testing.T) {
	docs.Description("validation should catch a port configuration that is out of range")
	tstValidatePort(t, "65536", "value '65536' is not a valid port number")
}

func TestValidateServerConfiguration_privileged(t *testing.T) {
	docs.Description("validation should not allow privileged ports")
	tstValidatePort(t, "1023", "value '1023' must be a nonprivileged port")
}

func createValidIdentityProviderConfiguration() identityProviderConfig {
	return identityProviderConfig{
		AuthorizationEndpoint: "https://example.com/auth",
		TokenEndpoint:         "https://example.com/token",
		EndSessionEndpoint:    "https://example.com/logout",
		TokenRequestTimeout:   time.Minute,
		AuthRequestTimeout:    time.Minute,
	}
}

func TestValidateIdentityProviderConfiguration_emptyAuthorizationEndpoint(t *testing.T) {
	docs.Description("validation should catch a missing authorization endpoint in identity provider config")
	errs := validationErrors{}
	config := createValidIdentityProviderConfiguration()
	config.AuthorizationEndpoint = ""
	validateIdentityProviderConfiguration(errs, config)
	require.Equal(t, 1, len(errs))
	require.Equal(t, []string{"value '' cannot not be empty"}, errs["identity_provider.authorization_endpoint"])
}

func TestValidateIdentityProviderConfiguration_emptyTokenEndpoint(t *testing.T) {
	docs.Description("validation should catch a missing token endpoint in identity provider config")
	errs := validationErrors{}
	config := createValidIdentityProviderConfiguration()
	config.TokenEndpoint = ""
	validateIdentityProviderConfiguration(errs, config)
	require.Equal(t, 1, len(errs))
	require.Equal(t, []string{"value '' cannot not be empty"}, errs["identity_provider.token_endpoint"])
}

func TestValidateIdentityProviderConfiguration_emptyEndSessionEndpoint(t *testing.T) {
	docs.Description("validation should catch a missing authorization endpoint in identity provider config")
	errs := validationErrors{}
	config := createValidIdentityProviderConfiguration()
	config.EndSessionEndpoint = ""
	validateIdentityProviderConfiguration(errs, config)
	require.Equal(t, 1, len(errs))
	require.Equal(t, []string{"value '' cannot not be empty"}, errs["identity_provider.end_session_endpoint"])
}

func TestValidateIdentityProviderConfiguration_zeroTokenRequestTimeout(t *testing.T) {
	docs.Description("validation should accept a zero token request timeout in identity provider config (means no timeout, not recommended)")
	errs := validationErrors{}
	config := createValidIdentityProviderConfiguration()
	config.TokenRequestTimeout = 0
	validateIdentityProviderConfiguration(errs, config)
	require.Equal(t, 0, len(errs))
}
func TestValidateIdentityProviderConfiguration_negativeTokenRequestTimeout(t *testing.T) {
	docs.Description("validation should catch a negative token request timeout in identity provider config")
	errs := validationErrors{}
	config := createValidIdentityProviderConfiguration()
	config.TokenRequestTimeout = -time.Second
	validateIdentityProviderConfiguration(errs, config)
	require.Equal(t, 1, len(errs))
	require.Equal(t, []string{"value '-1s' cannot be negative"}, errs["identity_provider.token_request_timeout"])
}

func TestValidateIdentityProviderConfiguration_zeroAuthRequestTimeout(t *testing.T) {
	docs.Description("validation should accept a zero auth request timeout in identity provider config")
	errs := validationErrors{}
	config := createValidIdentityProviderConfiguration()
	config.AuthRequestTimeout = 0
	validateIdentityProviderConfiguration(errs, config)
	require.Equal(t, 0, len(errs))
}

func TestValidateIdentityProviderConfiguration_negativeAuthRequestTimeout(t *testing.T) {
	docs.Description("validation should catch a negative auth request timeout in identity provider config")
	errs := validationErrors{}
	config := createValidIdentityProviderConfiguration()
	config.AuthRequestTimeout = -time.Second
	validateIdentityProviderConfiguration(errs, config)
	require.Equal(t, 1, len(errs))
	require.Equal(t, []string{"value '-1s' cannot be negative"}, errs["identity_provider.auth_request_timeout"])
}

func createValidApplicationConfig() ApplicationConfig {
	return ApplicationConfig{
		DisplayName:       "Test Application",
		Scope:             "test-scope",
		ClientId:          "test-client-id",
		ClientSecret:      "test-client-secret",
		DefaultDropoffUrl: "https://target.example.com/app",
		CookieName:        "ACookie",
		CookieDomain:      "example.com",
		CookiePath:        "/",
		CookieExpiry:      4 * time.Hour,
	}
}

func TestValidateApplicationConfigs_validSingle(t *testing.T) {
	docs.Description("validation should accept one valid application config")
	errs := validationErrors{}
	configs := map[string]ApplicationConfig{"test-application-config": createValidApplicationConfig()}
	validateApplicationConfigurations(errs, configs)
	require.Equal(t, 0, len(errs))
}

func TestValidateApplicationConfigs_validMultiple(t *testing.T) {
	docs.Description("validation should accept multiple valid application configs")
	errs := validationErrors{}
	configs := map[string]ApplicationConfig{"test-application-config-1": createValidApplicationConfig(), "test-application-config-2": createValidApplicationConfig()}
	validateApplicationConfigurations(errs, configs)
	require.Equal(t, 0, len(errs))
}

func TestValidateApplicationConfigs_empty(t *testing.T) {
	docs.Description("validation should require at least one application config")
	errs := validationErrors{}
	configs := map[string]ApplicationConfig{}
	validateApplicationConfigurations(errs, configs)
	require.Equal(t, 1, len(errs))
	require.Equal(t, []string{"value 'map[]' must contain at least one entry"}, errs["application_configs"])
}

func TestValidateApplicationConfigs_emptyDisplayName(t *testing.T) {
	docs.Description("validation should catch a missing application config display name")
	errs := validationErrors{}
	config := createValidApplicationConfig()
	config.DisplayName = ""
	configs := map[string]ApplicationConfig{"test-application-config": config}
	validateApplicationConfigurations(errs, configs)
	require.Equal(t, 1, len(errs))
	require.Equal(t, []string{"value '' cannot not be empty"}, errs["application_configs.test-application-config.display_name"])
}

func TestValidateApplicationConfigs_emptyScope(t *testing.T) {
	docs.Description("validation should catch a missing scope in application config")
	errs := validationErrors{}
	config := createValidApplicationConfig()
	config.Scope = ""
	configs := map[string]ApplicationConfig{"test-application-config": config}
	validateApplicationConfigurations(errs, configs)
	require.Equal(t, 1, len(errs))
	require.Equal(t, []string{"value '' cannot not be empty"}, errs["application_configs.test-application-config.scope"])
}

func TestValidateApplicationConfigs_emptyClientId(t *testing.T) {
	docs.Description("validation should catch a missing client ID in application config")
	errs := validationErrors{}
	config := createValidApplicationConfig()
	config.ClientId = ""
	configs := map[string]ApplicationConfig{"test-application-config": config}
	validateApplicationConfigurations(errs, configs)
	require.Equal(t, 1, len(errs))
	require.Equal(t, []string{"value '' cannot not be empty"}, errs["application_configs.test-application-config.client_id"])
}

func TestValidateApplicationConfigs_emptyClientSecret(t *testing.T) {
	docs.Description("validation should catch a missing client secret in application config")
	errs := validationErrors{}
	config := createValidApplicationConfig()
	config.ClientSecret = ""
	configs := map[string]ApplicationConfig{"test-application-config": config}
	validateApplicationConfigurations(errs, configs)
	require.Equal(t, 1, len(errs))
	require.Equal(t, []string{"value '' cannot not be empty"}, errs["application_configs.test-application-config.client_secret"])
}

func TestValidateApplicationConfigs_emptyDefaultDropoffUrl(t *testing.T) {
	docs.Description("validation should catch a missing default redirect URL in application config")
	errs := validationErrors{}
	config := createValidApplicationConfig()
	config.DefaultDropoffUrl = ""
	configs := map[string]ApplicationConfig{"test-application-config": config}
	validateApplicationConfigurations(errs, configs)
	require.Equal(t, 1, len(errs))
	require.Equal(t, []string{"value '' cannot not be empty"}, errs["application_configs.test-application-config.default_redirect_url"])
}

func TestValidateApplicationConfigs_emptyDropoffUrlPattern(t *testing.T) {
	docs.Description("validation should accept empty redirect URL pattern in application config")
	errs := validationErrors{}
	config := createValidApplicationConfig()
	config.DropoffUrlPattern = ""
	configs := map[string]ApplicationConfig{"test-application-config": config}
	validateApplicationConfigurations(errs, configs)
	require.Equal(t, 0, len(errs))
}

func TestValidateApplicationConfigs_validDropoffUrlPattern(t *testing.T) {
	docs.Description("validation should accept valid redirect URL pattern in application config")
	errs := validationErrors{}
	config := createValidApplicationConfig()
	config.DropoffUrlPattern = "https://reg.eurofurence.example.com/room/(\\?(foo=[a-z]+|bar=[0-9]{3,8}|&)+)?"
	configs := map[string]ApplicationConfig{"test-application-config": config}
	validateApplicationConfigurations(errs, configs)
	require.Equal(t, 0, len(errs))
}

func TestValidateApplicationConfigs_invalidDropoffUrlPattern(t *testing.T) {
	docs.Description("validation should catch an invalid redirect URL pattern in application config")
	errs := validationErrors{}
	config := createValidApplicationConfig()
	config.DropoffUrlPattern = "(iammissingaroundbracketattheendohno"
	configs := map[string]ApplicationConfig{"test-application-config": config}
	validateApplicationConfigurations(errs, configs)
	require.Equal(t, 1, len(errs))
	require.Equal(t, []string{"value '(iammissingaroundbracketattheendohno' must be a valid regular expression, but encountered compile error: error parsing regexp: missing closing ): `(iammissingaroundbracketattheendohno`)"}, errs["application_configs.test-application-config.redirect_url_pattern"])
}

func TestValidateApplicationConfigs_emptyCookieName(t *testing.T) {
	docs.Description("validation should catch a missing cookie name in application config")
	errs := validationErrors{}
	config := createValidApplicationConfig()
	config.CookieName = ""
	configs := map[string]ApplicationConfig{"test-application-config": config}
	validateApplicationConfigurations(errs, configs)
	require.Equal(t, 1, len(errs))
	require.Equal(t, []string{"value '' cannot not be empty"}, errs["application_configs.test-application-config.cookie_name"])
}

func TestValidateApplicationConfigs_emptyCookieDomain(t *testing.T) {
	docs.Description("validation should catch a missing cookie domain in application config")
	errs := validationErrors{}
	config := createValidApplicationConfig()
	config.CookieDomain = ""
	configs := map[string]ApplicationConfig{"test-application-config": config}
	validateApplicationConfigurations(errs, configs)
	require.Equal(t, 1, len(errs))
	require.Equal(t, []string{"value '' cannot not be empty"}, errs["application_configs.test-application-config.cookie_domain"])
}

func TestValidateApplicationConfigs_emptyCookiePath(t *testing.T) {
	docs.Description("validation should catch a missing cookie path in application config")
	errs := validationErrors{}
	config := createValidApplicationConfig()
	config.CookiePath = ""
	configs := map[string]ApplicationConfig{"test-application-config": config}
	validateApplicationConfigurations(errs, configs)
	require.Equal(t, 1, len(errs))
	require.Equal(t, []string{"value '' cannot not be empty, use '/' for all paths"}, errs["application_configs.test-application-config.cookie_path"])
}

func TestValidateApplicationConfigs_negativeCookieExpiry(t *testing.T) {
	docs.Description("validation should catch a negative cookie expiry in application config")
	errs := validationErrors{}
	config := createValidApplicationConfig()
	config.CookieExpiry = -1 * time.Hour
	configs := map[string]ApplicationConfig{"test-application-config": config}
	validateApplicationConfigurations(errs, configs)
	require.Equal(t, 1, len(errs))
	require.Equal(t, []string{"value '-1h0m0s' must be positive, try '1h' or '5m'"}, errs["application_configs.test-application-config.cookie_expiry"])
}

func TestValidateApplicationConfigs_zeroCookieExpiry(t *testing.T) {
	docs.Description("validation should catch a zero cookie expiry in application config")
	errs := validationErrors{}
	config := createValidApplicationConfig()
	config.CookieExpiry = 0
	configs := map[string]ApplicationConfig{"test-application-config": config}
	validateApplicationConfigurations(errs, configs)
	require.Equal(t, 1, len(errs))
	require.Equal(t, []string{"value '0s' must be positive, try '1h' or '5m'"}, errs["application_configs.test-application-config.cookie_expiry"])
}
