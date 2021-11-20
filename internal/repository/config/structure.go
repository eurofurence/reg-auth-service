package config

type conf struct {
	Server             serverConfig                 `yaml:"server"`
	Security           securityConfig               `yaml:"security"`
	ApplicationConfigs map[string]applicationConfig `yaml:"application_configs"`
}

type serverConfig struct {
	Port string `yaml:"port"`
}

type securityConfig struct {
	DisableCors bool `yaml:"disable_cors"`
}

type applicationConfig struct {
	DisplayName           string `yaml:"display_name"`
	AuthorizationEndpoint string `yaml:"authorization_endpoint"`
	Scope                 string `yaml:"scope"`
	ClientId              string `yaml:"client_id"`
	ClientSecret          string `yaml:"client_secret"`
	DefaultRedirectUrl    string `yaml:"default_redirect_url"`
	CodeChallengeMethod   string `yaml:"code_challenge_method"`
}

type validationErrors map[string][]string
