package config

import (
	"errors"
	"flag"
	aulogging "github.com/StephanHCB/go-autumn-logging"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"net/url"
	"sort"
)

var (
	configurationData     *conf
	configurationFilename string
	ecsLogging            bool
)

var (
	ErrorConfigArgumentMissing = errors.New("configuration file argument missing. Please specify using -config argument. Aborting")
	ErrorConfigFile            = errors.New("failed to read or parse configuration file. Aborting")
)

func init() {
	configurationData = &conf{}

	flag.StringVar(&configurationFilename, "config", "", "config file path")
	flag.BoolVar(&ecsLogging, "ecs-json-logging", false, "switch to structured json logging")
}

// ParseCommandLineFlags is exposed separately so you can skip it for tests
func ParseCommandLineFlags() {
	flag.Parse()
}

func logValidationErrors(errs url.Values) error {
	if len(errs) != 0 {
		var keys []string
		for key := range errs {
			keys = append(keys, key)
		}
		sort.Strings(keys)

		for _, k := range keys {
			key := k
			val := errs[k]
			for _, errorvalue := range val {
				aulogging.Logger.NoCtx().Error().Printf("configuration error: %s: %s", key, errorvalue)
			}
		}
		return errors.New("configuration validation error, see log output for details")
	}

	return nil
}

func configuration() *conf {
	return configurationData
}

func setConfigurationDefaults(c *conf) {
	if c.Server.Port == "" {
		c.Server.Port = "8081"
	}
	if c.Server.ReadTimeout <= 0 {
		c.Server.ReadTimeout = 5
	}
	if c.Server.WriteTimeout <= 0 {
		c.Server.WriteTimeout = 5
	}
	if c.Server.IdleTimeout <= 0 {
		c.Server.IdleTimeout = 5
	}
	if c.Logging.Severity == "" {
		c.Logging.Severity = "INFO"
	}
}

func validateConfiguration(newConfigurationData *conf) error {
	errs := url.Values{}

	validateServerConfiguration(errs, newConfigurationData.Server)
	validateLoggingConfiguration(errs, newConfigurationData.Logging)
	validateSecurityConfiguration(errs, newConfigurationData.Security)
	validateDropoffEndpointUrl(errs, newConfigurationData.DropoffEndpointUrl)
	validateIdentityProviderConfiguration(errs, newConfigurationData.IdentityProvider)
	validateApplicationConfigurations(errs, newConfigurationData.ApplicationConfigs)

	return logValidationErrors(errs)
}

func ParseAndOverwriteConfig(yamlFile []byte) error {
	newConfigurationData := &conf{}
	err := yaml.UnmarshalStrict(yamlFile, newConfigurationData)
	if err != nil {
		return err
	}

	setConfigurationDefaults(newConfigurationData)

	err = validateConfiguration(newConfigurationData)
	if err != nil {
		return err
	}

	configurationData = newConfigurationData
	return nil
}

// exposed for testing

func LoadConfiguration(filename string) error {
	yamlFile, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}

	err = ParseAndOverwriteConfig(yamlFile)
	return err
}

func StartupLoadConfiguration() error {
	aulogging.Logger.NoCtx().Info().Print("Reading configuration...")
	if configurationFilename == "" {
		// cannot use logging package here as this would create a circular dependency (logging needs config)
		aulogging.Logger.NoCtx().Error().Print("Configuration file argument missing. Please specify using -config argument. Aborting.")
		return ErrorConfigArgumentMissing
	}
	err := LoadConfiguration(configurationFilename)
	if err != nil {
		// cannot use logging package here as this would create a circular dependency (logging needs config)
		aulogging.Logger.NoCtx().Error().Print("Error reading or parsing configuration file. Aborting.")
		return ErrorConfigFile
	}
	return nil
}
