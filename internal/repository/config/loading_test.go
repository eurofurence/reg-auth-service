package config

import (
	"testing"

	"github.com/eurofurence/reg-auth-service/docs"
	"github.com/stretchr/testify/require"
)

func TestLoadConfiguration_noFilename(t *testing.T) {
	docs.Description("empty configuration filename is an error")
	err := StartupLoadConfiguration()
	require.NotNil(t, err)
	require.Equal(t, "configuration file argument missing. Please specify using -config argument. Aborting", err.Error())
}
