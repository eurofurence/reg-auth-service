package config

type conf struct {
	Server             serverConfig                 `yaml:"server"`
	Security           securityConfig               `yaml:"security"`
	IdentityProvider   identityProviderConfig       `yaml:"identity_provider"`
	ApplicationConfigs map[string]ApplicationConfig `yaml:"application_configs"`
}

type serverConfig struct {
	Port string `yaml:"port"`
}

type securityConfig struct {
	DisableCors bool `yaml:"disable_cors"`
}

type identityProviderConfig struct {
	AuthorizationEndpoint   string `yaml:"authorization_endpoint"`
	TokenEndpoint           string `yaml:"token_endpoint"`
	EndSessionEndpoint      string `yaml:"end_session_endpoint"`
	CircuitBreakerTimeoutMS int    `yaml:"circuit_breaker_timeout_ms"`
	AuthRequestTimeoutS     int    `yaml:"auth_request_timeout_s"`
}

type ApplicationConfig struct {
	DisplayName         string `yaml:"display_name"`
	Scope               string `yaml:"scope"`
	ClientId            string `yaml:"client_id"`
	ClientSecret        string `yaml:"client_secret"`
	DefaultRedirectUrl  string `yaml:"default_redirect_url"`
	RedirectUrlPattern  string `yaml:"redirect_url_pattern"`
}

type validationErrors map[string][]string