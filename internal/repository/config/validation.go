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
	if ipc.CircuitBreakerTimeoutMS <= 0 {
		addError(errs, "identity_provider.circuit_breaker_timeout_ms", ipc.CircuitBreakerTimeoutMS, "must be greater than 0")
	}
	if ipc.AuthRequestTimeoutS <= 0 {
		addError(errs, "identity_provider.auth_request_timeout_s", ipc.AuthRequestTimeoutS, "must be greater than 0")
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
		if ac.DefaultRedirectUrl == "" {
			addError(errs, fmt.Sprintf("application_configs.%s.default_redirect_url", name), ac.DefaultRedirectUrl, "cannot not be empty")
		}
		if ac.RedirectUrlPattern != "" {
			if _, regexpError := regexp.Compile(strings.ReplaceAll(ac.RedirectUrlPattern, "/", "\\/")); regexpError != nil {
				addError(errs, fmt.Sprintf("application_configs.%s.redirect_url_pattern", name), ac.RedirectUrlPattern, fmt.Sprintf("must be a valid regular expression, but encountered compile error: %s)", regexpError))
			}
		}
	}
}
