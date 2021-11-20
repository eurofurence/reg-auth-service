package config

func ServerAddr() string {
	return ":" + configuration().Server.Port
}

func IsCorsDisabled() bool {
	return configuration().Security.DisableCors
}
