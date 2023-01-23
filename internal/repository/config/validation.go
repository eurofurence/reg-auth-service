package config

import (
	"crypto/rsa"
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"net/url"
	"regexp"
	"strconv"
	"strings"
)

func addError(errs url.Values, key string, value interface{}, message string) {
	errs[key] = append(errs[key], fmt.Sprintf("value '%v' %s", value, message))
}

func notInAllowedValues(allowed []string, value string) bool {
	for _, v := range allowed {
		if v == value {
			return false
		}
	}
	return true
}

func checkIntValueRange(errs *url.Values, min int, max int, key string, value int) {
	if value < min || value > max {
		errs.Add(key, fmt.Sprintf("%s field must be an integer at least %d and at most %d", key, min, max))
	}
}

func validateDropoffEndpointUrl(errs url.Values, value string) {
	if value == "" {
		addError(errs, "dropoff_endpoint_url", value, "cannot not be empty")
	}
}

func validateServerConfiguration(errs url.Values, sc ServerConfig) {
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
	checkIntValueRange(&errs, 1, 300, "server.read_timeout_seconds", sc.ReadTimeout)
	checkIntValueRange(&errs, 1, 300, "server.write_timeout_seconds", sc.WriteTimeout)
	checkIntValueRange(&errs, 1, 300, "server.idle_timeout_seconds", sc.IdleTimeout)
}

func validateSecurityConfiguration(errs url.Values, c SecurityConfig) {
	parsedKeySet = make([]*rsa.PublicKey, 0)
	for i, keyStr := range c.Oidc.TokenPublicKeysPEM {
		publicKeyPtr, err := jwt.ParseRSAPublicKeyFromPEM([]byte(keyStr))
		if err != nil {
			errs.Add(fmt.Sprintf("security.oidc.token_public_keys_PEM[%d]", i), fmt.Sprintf("failed to parse RSA public key in PEM format: %s", err.Error()))
		} else {
			parsedKeySet = append(parsedKeySet, publicKeyPtr)
		}
	}

	if c.Cors.DisableCors && c.Cors.InsecureCookies {
		errs.Add("security.cors.disable", "not compatible with security.cors.insecure_cookies, because SameSitePolicy None only works with secure cookies")
	}
}

var allowedSeverities = []string{"DEBUG", "INFO", "WARN", "ERROR"}

func validateLoggingConfiguration(errs url.Values, c LoggingConfig) {
	if notInAllowedValues(allowedSeverities[:], c.Severity) {
		errs.Add("logging.severity", "must be one of DEBUG, INFO, WARN, ERROR")
	}
}

func validateIdentityProviderConfiguration(errs url.Values, ipc IdentityProviderConfig) {
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

func validateApplicationConfigurations(errs url.Values, acs map[string]ApplicationConfig) {
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
