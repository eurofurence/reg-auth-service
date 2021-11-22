package config

import (
	"fmt"
	"time"
)

func ServerAddr() string {
	return ":" + configuration().Server.Port
}

func IsCorsDisabled() bool {
	return configuration().Security.DisableCors
}

func TokenEndpoint() string {
	return configuration().IdentityProvider.TokenEndpoint
}

func AuthorizationEndpoint() string {
	return configuration().IdentityProvider.AuthorizationEndpoint
}

func TokenRequestTimeout() time.Duration {
	return configuration().IdentityProvider.TokenRequestTimeout
}

func AuthRequestTimeout() time.Duration {
	return configuration().IdentityProvider.AuthRequestTimeout
}

func GetApplicationConfig(applicationName string) (ApplicationConfig, error) {
	appConfig, found := configuration().ApplicationConfigs[applicationName]
	if found {
		return appConfig, nil
	} else {
		return ApplicationConfig{}, fmt.Errorf("no application configured for applicationName %s", applicationName)
	}
}
