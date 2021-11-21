package config

import "fmt"

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

func CircuitBreakerTimeoutMilliseconds() int {
	return configuration().IdentityProvider.CircuitBreakerTimeoutMS
}

func GetApplicationConfig(applicationName string) (ApplicationConfig, error) {
	appConfig, found := configuration().ApplicationConfigs[applicationName]
	if found {
		return appConfig, nil
	} else {
		return ApplicationConfig{}, fmt.Errorf("no application configured for applicationName %s", applicationName)
	}
}
