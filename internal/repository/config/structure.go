package config

import "time"

type conf struct {
	Server             serverConfig                 `yaml:"server"`
	Logging            loggingConfig                `yaml:"logging"`
	Security           securityConfig               `yaml:"security"`
	DropoffEndpointUrl string                       `yaml:"dropoff_endpoint_url"`
	IdentityProvider   identityProviderConfig       `yaml:"identity_provider"`
	ApplicationConfigs map[string]ApplicationConfig `yaml:"application_configs"`
}

type serverConfig struct {
	Address      string `yaml:"address"`
	Port         string `yaml:"port"`
	ReadTimeout  int    `yaml:"read_timeout_seconds"`
	WriteTimeout int    `yaml:"write_timeout_seconds"`
	IdleTimeout  int    `yaml:"idle_timeout_seconds"`
}

type securityConfig struct {
	DisableCors bool `yaml:"disable_cors"`
}

type loggingConfig struct {
	Severity string `yaml:"severity"`
}

type identityProviderConfig struct {
	AuthorizationEndpoint string        `yaml:"authorization_endpoint"`
	TokenEndpoint         string        `yaml:"token_endpoint"`
	EndSessionEndpoint    string        `yaml:"end_session_endpoint"`
	TokenRequestTimeout   time.Duration `yaml:"token_request_timeout"`
	AuthRequestTimeout    time.Duration `yaml:"auth_request_timeout"`
}

type ApplicationConfig struct {
	DisplayName       string        `yaml:"display_name"`
	Scope             string        `yaml:"scope"`
	ClientId          string        `yaml:"client_id"`
	ClientSecret      string        `yaml:"client_secret"`
	DefaultDropoffUrl string        `yaml:"default_dropoff_url"`
	DropoffUrlPattern string        `yaml:"dropoff_url_pattern"`
	CookieName        string        `yaml:"cookie_name"`
	CookieDomain      string        `yaml:"cookie_domain"`
	CookiePath        string        `yaml:"cookie_path"`
	CookieExpiry      time.Duration `yaml:"cookie_expiry"`
}
