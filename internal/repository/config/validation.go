package config

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

func addError(errs validationErrors, key string, value interface{}, message string) {
	errs[key] = append(errs[key], fmt.Sprintf("value '%v' %s", value, message))
}

func validateDropoffEndpointUrl(errs validationErrors, value string) {
	if value == "" {
		addError(errs, "dropoff_endpoint_url", value, "cannot not be empty")
	}
}

func validateServerConfiguration(errs validationErrors, sc serverConfig) {
	if sc.Port == "" {
		addError(errs, "server.port", sc.Port, "cannot be empty")
	} else {
		port, err := strconv.ParseUint(sc.Port, 10, 16)
		if err != nil {
			addError(errs, "server.port", sc.Port, "is not a valid port number")
		} else if port <= 1024 {
			addError(errs, "server.port", sc.Port, "must be a nonprivileged port")
		}
	}
}

func validateSecurityConfiguration(errs validationErrors, sc securityConfig) {
}

func validateIdentityProviderConfiguration(errs validationErrors, ipc identityProviderConfig) {
	if ipc.AuthorizationEndpoint == "" {
		addError(errs, "identity_provider.authorization_endpoint", ipc.AuthorizationEndpoint, "cannot not be empty")
	}
	if ipc.TokenEndpoint == "" {
		addError(errs, "identity_provider.token_endpoint", ipc.TokenEndpoint, "cannot not be empty")
	}
	if ipc.EndSessionEndpoint == "" {
		addError(errs, "identity_provider.end_session_endpoint", ipc.EndSessionEndpoint, "cannot not be empty")
	}
	if ipc.TokenRequestTimeout < 0 {
		addError(errs, "identity_provider.token_request_timeout", ipc.TokenRequestTimeout, "cannot be negative")
	}
	if ipc.AuthRequestTimeout < 0 {
		addError(errs, "identity_provider.auth_request_timeout", ipc.AuthRequestTimeout, "cannot be negative")
	}
}

func validateApplicationConfigurations(errs validationErrors, acs map[string]ApplicationConfig) {
	if len(acs) == 0 {
		addError(errs, "application_configs", acs, "must contain at least one entry")
	}
	for name, ac := range acs {
		if ac.DisplayName == "" {
			addError(errs, fmt.Sprintf("application_configs.%s.display_name", name), ac.DisplayName, "cannot not be empty")
		}
		if ac.Scope == "" {
			addError(errs, fmt.Sprintf("application_configs.%s.scope", name), ac.Scope, "cannot not be empty")
		}
		if ac.ClientId == "" {
			addError(errs, fmt.Sprintf("application_configs.%s.client_id", name), ac.ClientId, "cannot not be empty")
		}
		if ac.ClientSecret == "" {
			addError(errs, fmt.Sprintf("application_configs.%s.client_secret", name), ac.ClientSecret, "cannot not be empty")
		}
		if ac.DefaultDropoffUrl == "" {
			addError(errs, fmt.Sprintf("application_configs.%s.default_redirect_url", name), ac.DefaultDropoffUrl, "cannot not be empty")
		}
		if ac.DropoffUrlPattern != "" {
			if _, regexpError := regexp.Compile(strings.ReplaceAll(ac.DropoffUrlPattern, "/", "\\/")); regexpError != nil {
				addError(errs, fmt.Sprintf("application_configs.%s.redirect_url_pattern", name), ac.DropoffUrlPattern, fmt.Sprintf("must be a valid regular expression, but encountered compile error: %s)", regexpError))
			}
		}
		if ac.CookieName == "" {
			addError(errs, fmt.Sprintf("application_configs.%s.cookie_name", name), ac.CookieName, "cannot not be empty")
		}
		if ac.CookieDomain == "" {
			addError(errs, fmt.Sprintf("application_configs.%s.cookie_domain", name), ac.CookieDomain, "cannot not be empty")
		}
		if ac.CookiePath == "" {
			addError(errs, fmt.Sprintf("application_configs.%s.cookie_path", name), ac.CookiePath, "cannot not be empty, use '/' for all paths")
		}
		if ac.CookieExpiry <= 0 {
			addError(errs, fmt.Sprintf("application_configs.%s.cookie_expiry", name), ac.CookieExpiry, "must be positive, try '1h' or '5m'")
		}
	}
}
