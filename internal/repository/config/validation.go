package config

import (
	"fmt"
	"strconv"
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

func validateApplicationConfigurations(errs validationErrors, acs []applicationConfig) {
	if len(acs) == 0 {
		addError(errs, "application_configs", acs, "must contain at least one entry")
	}
	for i, ac := range acs {
		if ac.Name == "" {
			addError(errs, fmt.Sprintf("application_configs[%d].name", i), ac.Name, "cannot not be empty")
		} else {
			//TODO: check uniqueness of application config names
		}
		if ac.AuthorizationEndpoint == "" {
			addError(errs, fmt.Sprintf("application_configs[%d].authorization_endpoint", i), ac.AuthorizationEndpoint, "cannot not be empty")
		}
		if ac.Scope == "" {
			addError(errs, fmt.Sprintf("application_configs[%d].scope", i), ac.Scope, "cannot not be empty")
		}
		if ac.ClientId == "" {
			addError(errs, fmt.Sprintf("application_configs[%d].client_id", i), ac.ClientId, "cannot not be empty")
		}
		if ac.ClientSecret == "" {
			addError(errs, fmt.Sprintf("application_configs[%d].client_secret", i), ac.ClientSecret, "cannot not be empty")
		}
		if ac.RedirectUrl == "" {
			addError(errs, fmt.Sprintf("application_configs[%d].redirect_url", i), ac.RedirectUrl, "cannot not be empty")
		}
		if ac.CodeChallengeMethod != "S256" && ac.CodeChallengeMethod != "" {
			addError(errs, fmt.Sprintf("application_configs[%d].code_challenge_method", i), ac.CodeChallengeMethod, "must be empty or S256")
		}
	}
}
