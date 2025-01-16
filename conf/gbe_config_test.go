package conf

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetConfig(t *testing.T) {
	tempDir := t.TempDir()

	SetConfFilePath(tempDir)
	SetConfFileName("test_config.json")

	configData, err := json.Marshal(testConfig)
	require.NoError(t, err)

	configFilePath := tempDir + "/test_config.json"
	err = os.WriteFile(configFilePath, configData, 0644)
	require.NoError(t, err)

	loadedConfig := GetConfig()

	assert.Equal(t, testConfig, *loadedConfig, "The loaded configuration does not match the expected configuration")
}
