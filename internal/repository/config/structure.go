package config

import "time"

type (
	// Application is the root configuration type
	Application struct {
		Service            ServiceConfig                `yaml:"service"`
		Server             ServerConfig                 `yaml:"server"`
		Security           SecurityConfig               `yaml:"security"`
		Logging            LoggingConfig                `yaml:"logging"`
		IdentityProvider   IdentityProviderConfig       `yaml:"identity_provider"`
		ApplicationConfigs map[string]ApplicationConfig `yaml:"application_configs"`
	}

	// ServiceConfig contains configuration values
	// for service related tasks. E.g. URLs to downstream services
	ServiceConfig struct {
		Name               string `yaml:"name"`
		DropoffEndpointUrl string `yaml:"dropoff_endpoint_url"` // externally visible url to my "dropoff" endpoint
	}

	// ServerConfig contains all values for http configuration
	ServerConfig struct {
		Address      string `yaml:"address"`
		Port         string `yaml:"port"`
		ReadTimeout  int    `yaml:"read_timeout_seconds"`
		WriteTimeout int    `yaml:"write_timeout_seconds"`
		IdleTimeout  int    `yaml:"idle_timeout_seconds"`
	}

	// SecurityConfig configures everything related to security
	SecurityConfig struct {
		Cors CorsConfig `yaml:"cors"`
	}

	CorsConfig struct {
		DisableCors bool   `yaml:"disable"`
		AllowOrigin string `yaml:"allow_origin"`
	}

	// LoggingConfig configures logging
	LoggingConfig struct {
		Severity string `yaml:"severity"`
	}

	// IdentityProviderConfig provides information about an OpenID Connect identity provider
	IdentityProviderConfig struct {
		AuthorizationEndpoint string        `yaml:"authorization_endpoint"`
		TokenEndpoint         string        `yaml:"token_endpoint"`
		EndSessionEndpoint    string        `yaml:"end_session_endpoint"`
		UserInfoEndpoint      string        `yaml:"user_info_endpoint"`
		KeySetEndpoint        string        `yaml:"key_set_endpoint"`
		TokenRequestTimeout   time.Duration `yaml:"token_request_timeout"`
		AuthRequestTimeout    time.Duration `yaml:"auth_request_timeout"`
	}

	// ApplicationConfig configures an OpenID Connect client.
	ApplicationConfig struct {
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
)
