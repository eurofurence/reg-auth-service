package config

import (
	"testing"

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
		AuthorizationEndpoint:   "https://example.com/auth",
		TokenEndpoint:           "https://example.com/token",
		EndSessionEndpoint:      "https://example.com/logout",
		CircuitBreakerTimeoutMS: 21,
		AuthRequestTimeoutS:     42,
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

func TestValidateIdentityProviderConfiguration_zeroCircuitBreakerTimeout(t *testing.T) {
	docs.Description("validation should catch a zero circuit breaker timeout in identity provider config")
	errs := validationErrors{}
	config := createValidIdentityProviderConfiguration()
	config.CircuitBreakerTimeoutMS = 0
	validateIdentityProviderConfiguration(errs, config)
	require.Equal(t, 1, len(errs))
	require.Equal(t, []string{"value '0' must be greater than 0"}, errs["identity_provider.circuit_breaker_timeout_ms"])
}
func TestValidateIdentityProviderConfiguration_negativeCircuitBreakerTimeout(t *testing.T) {
	docs.Description("validation should catch a negative circuit breaker timeout in identity provider config")
	errs := validationErrors{}
	config := createValidIdentityProviderConfiguration()
	config.CircuitBreakerTimeoutMS = -21
	validateIdentityProviderConfiguration(errs, config)
	require.Equal(t, 1, len(errs))
	require.Equal(t, []string{"value '-21' must be greater than 0"}, errs["identity_provider.circuit_breaker_timeout_ms"])
}

func TestValidateIdentityProviderConfiguration_zeroAuthRequestTimeout(t *testing.T) {
	docs.Description("validation should catch a zero auth request timeout in identity provider config")
	errs := validationErrors{}
	config := createValidIdentityProviderConfiguration()
	config.AuthRequestTimeoutS = 0
	validateIdentityProviderConfiguration(errs, config)
	require.Equal(t, 1, len(errs))
	require.Equal(t, []string{"value '0' must be greater than 0"}, errs["identity_provider.auth_request_timeout_s"])
}
func TestValidateIdentityProviderConfiguration_negativeAuthRequestTimeout(t *testing.T) {
	docs.Description("validation should catch a negative auth request timeout in identity provider config")
	errs := validationErrors{}
	config := createValidIdentityProviderConfiguration()
	config.AuthRequestTimeoutS = -21
	validateIdentityProviderConfiguration(errs, config)
	require.Equal(t, 1, len(errs))
	require.Equal(t, []string{"value '-21' must be greater than 0"}, errs["identity_provider.auth_request_timeout_s"])
}

func createValidApplicationConfig() ApplicationConfig {
	return ApplicationConfig{
		DisplayName:         "Test Application",
		Scope:               "test-scope",
		ClientId:            "test-client-id",
		ClientSecret:        "test-client-secret",
		DefaultRedirectUrl:  "https://target.example.com/app",
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

func TestValidateApplicationConfigs_emptyDefaultRedirectUrl(t *testing.T) {
	docs.Description("validation should catch a missing default redirect URL in application config")
	errs := validationErrors{}
	config := createValidApplicationConfig()
	config.DefaultRedirectUrl = ""
	configs := map[string]ApplicationConfig{"test-application-config": config}
	validateApplicationConfigurations(errs, configs)
	require.Equal(t, 1, len(errs))
	require.Equal(t, []string{"value '' cannot not be empty"}, errs["application_configs.test-application-config.default_redirect_url"])
}

func TestValidateApplicationConfigs_emptyRedirectUrlPattern(t *testing.T) {
	docs.Description("validation should accept empty redirect URL pattern in application config")
	errs := validationErrors{}
	config := createValidApplicationConfig()
	config.RedirectUrlPattern = ""
	configs := map[string]ApplicationConfig{"test-application-config": config}
	validateApplicationConfigurations(errs, configs)
	require.Equal(t, 0, len(errs))
}

func TestValidateApplicationConfigs_validRedirectUrlPattern(t *testing.T) {
	docs.Description("validation should accept valid redirect URL pattern in application config")
	errs := validationErrors{}
	config := createValidApplicationConfig()
	config.RedirectUrlPattern = "https://reg.eurofurence.example.com/room/(\\?(foo=[a-z]+|bar=[0-9]{3,8}|&)+)?"
	configs := map[string]ApplicationConfig{"test-application-config": config}
	validateApplicationConfigurations(errs, configs)
	require.Equal(t, 0, len(errs))
}

func TestValidateApplicationConfigs_invalidRedirectUrlPattern(t *testing.T) {
	docs.Description("validation should catch an invalid redirect URL pattern in application config")
	errs := validationErrors{}
	config := createValidApplicationConfig()
	config.RedirectUrlPattern = "(iammissingaroundbracketattheendohno"
	configs := map[string]ApplicationConfig{"test-application-config": config}
	validateApplicationConfigurations(errs, configs)
	require.Equal(t, 1, len(errs))
	require.Equal(t, []string{"value '(iammissingaroundbracketattheendohno' must be a valid regular expression, but encountered compile error: error parsing regexp: missing closing ): `(iammissingaroundbracketattheendohno`)"}, errs["application_configs.test-application-config.redirect_url_pattern"])
}
