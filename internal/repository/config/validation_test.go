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

func createValidApplicationConfig() applicationConfig {
	return applicationConfig{
		DisplayName:           "Test Application",
		AuthorizationEndpoint: "https://idp.example.com/auth",
		Scope:                 "test-scope",
		ClientId:              "test-client-id",
		ClientSecret:          "test-client-secret",
		RedirectUrl:           "https://target.example.com/app",
		CodeChallengeMethod:   "S256",
	}
}

func TestValidateApplicationConfigs_validSingle(t *testing.T) {
	docs.Description("validation should accept one valid application config")
	errs := validationErrors{}
	configs := map[string]applicationConfig{"test-application-config": createValidApplicationConfig()}
	validateApplicationConfigurations(errs, configs)
	require.Equal(t, 0, len(errs))
}

func TestValidateApplicationConfigs_validMultiple(t *testing.T) {
	docs.Description("validation should accept multiple valid application configs")
	errs := validationErrors{}
	configs := map[string]applicationConfig{"test-application-config-1": createValidApplicationConfig(), "test-application-config-2": createValidApplicationConfig()}
	validateApplicationConfigurations(errs, configs)
	require.Equal(t, 0, len(errs))
}

func TestValidateApplicationConfigs_empty(t *testing.T) {
	docs.Description("validation should require at least one application config")
	errs := validationErrors{}
	configs := map[string]applicationConfig{}
	validateApplicationConfigurations(errs, configs)
	require.Equal(t, 1, len(errs))
	require.Equal(t, []string{"value 'map[]' must contain at least one entry"}, errs["application_configs"])
}

func TestValidateApplicationConfigs_emptyDisplayName(t *testing.T) {
	docs.Description("validation should catch a missing application config display name")
	errs := validationErrors{}
	config := createValidApplicationConfig()
	config.DisplayName = ""
	configs := map[string]applicationConfig{"test-application-config": config}
	validateApplicationConfigurations(errs, configs)
	require.Equal(t, 1, len(errs))
	require.Equal(t, []string{"value '' cannot not be empty"}, errs["application_configs.test-application-config.display_name"])
}

func TestValidateApplicationConfigs_emptyAuthorizationEndpoint(t *testing.T) {
	docs.Description("validation should catch a missing authorization endpoint in application config")
	errs := validationErrors{}
	config := createValidApplicationConfig()
	config.AuthorizationEndpoint = ""
	configs := map[string]applicationConfig{"test-application-config": config}
	validateApplicationConfigurations(errs, configs)
	require.Equal(t, 1, len(errs))
	require.Equal(t, []string{"value '' cannot not be empty"}, errs["application_configs.test-application-config.authorization_endpoint"])
}

func TestValidateApplicationConfigs_emptyScope(t *testing.T) {
	docs.Description("validation should catch a missing scope in application config")
	errs := validationErrors{}
	config := createValidApplicationConfig()
	config.Scope = ""
	configs := map[string]applicationConfig{"test-application-config": config}
	validateApplicationConfigurations(errs, configs)
	require.Equal(t, 1, len(errs))
	require.Equal(t, []string{"value '' cannot not be empty"}, errs["application_configs.test-application-config.scope"])
}

func TestValidateApplicationConfigs_emptyClientId(t *testing.T) {
	docs.Description("validation should catch a missing client ID in application config")
	errs := validationErrors{}
	config := createValidApplicationConfig()
	config.ClientId = ""
	configs := map[string]applicationConfig{"test-application-config": config}
	validateApplicationConfigurations(errs, configs)
	require.Equal(t, 1, len(errs))
	require.Equal(t, []string{"value '' cannot not be empty"}, errs["application_configs.test-application-config.client_id"])
}

func TestValidateApplicationConfigs_emptyClientSecret(t *testing.T) {
	docs.Description("validation should catch a missing client secret in application config")
	errs := validationErrors{}
	config := createValidApplicationConfig()
	config.ClientSecret = ""
	configs := map[string]applicationConfig{"test-application-config": config}
	validateApplicationConfigurations(errs, configs)
	require.Equal(t, 1, len(errs))
	require.Equal(t, []string{"value '' cannot not be empty"}, errs["application_configs.test-application-config.client_secret"])
}

func TestValidateApplicationConfigs_emptyRedirectUrl(t *testing.T) {
	docs.Description("validation should catch a missing redirect URL in application config")
	errs := validationErrors{}
	config := createValidApplicationConfig()
	config.RedirectUrl = ""
	configs := map[string]applicationConfig{"test-application-config": config}
	validateApplicationConfigurations(errs, configs)
	require.Equal(t, 1, len(errs))
	require.Equal(t, []string{"value '' cannot not be empty"}, errs["application_configs.test-application-config.redirect_url"])
}

func TestValidateApplicationConfigs_emptyCodeChallengeMethod(t *testing.T) {
	docs.Description("validation should accept an empty code challenge method in application config")
	errs := validationErrors{}
	config := createValidApplicationConfig()
	config.CodeChallengeMethod = ""
	configs := map[string]applicationConfig{"test-application-config": config}
	validateApplicationConfigurations(errs, configs)
	require.Equal(t, 0, len(errs))
}

func TestValidateApplicationConfigs_invalidCodeChallengeMethod(t *testing.T) {
	docs.Description("validation should catch an invalid code challenge method in application config")
	errs := validationErrors{}
	config := createValidApplicationConfig()
	config.CodeChallengeMethod = "INV"
	configs := map[string]applicationConfig{"invalid-test-application-config": config}
	validateApplicationConfigurations(errs, configs)
	require.Equal(t, 1, len(errs))
	require.Equal(t, []string{"value 'INV' must be empty or S256"}, errs["application_configs.invalid-test-application-config.code_challenge_method"])
}

func TestValidateApplicationConfigs_invalidCodeChallengeMethodWithMultipleConfigs(t *testing.T) {
	docs.Description("validation should catch an invalid code challenge method in other application config")
	errs := validationErrors{}
	config := createValidApplicationConfig()
	config.CodeChallengeMethod = "INV"
	configs := map[string]applicationConfig{"test-application-config": createValidApplicationConfig(), "invalid-test-application-config": config}
	validateApplicationConfigurations(errs, configs)
	require.Equal(t, 1, len(errs))
	require.Equal(t, []string{"value 'INV' must be empty or S256"}, errs["application_configs.invalid-test-application-config.code_challenge_method"])
}
