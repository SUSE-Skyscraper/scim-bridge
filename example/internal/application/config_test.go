package application

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewConfigurator(t *testing.T) {
	configDir, err := filepath.Abs("../../")
	assert.Nil(t, err)

	configurator := NewConfigurator(configDir)
	assert.NotNil(t, configurator)
}

func TestConfigurator_Parse(t *testing.T) {
	configDir, err := filepath.Abs("../../")
	assert.Nil(t, err)

	configurator := NewConfigurator(configDir)
	assert.NotNil(t, configurator)

	config, err := configurator.Parse()
	assert.Nil(t, err)
	assert.NotNil(t, config)
}

func TestConfigurator_Parse_Bad(t *testing.T) {
	configDir, err := filepath.Abs("../../../../testdata/bad-config")
	assert.Nil(t, err)

	configurator := NewConfigurator(configDir)
	assert.NotNil(t, configurator)

	config, err := configurator.Parse()
	assert.NotNil(t, err)
	assert.NotNil(t, config)
}

func TestConfigurator_Defaults(t *testing.T) {
	configDir, err := filepath.Abs("../../../../testdata/null-config")
	assert.Nil(t, err)

	configurator := NewConfigurator(configDir)
	assert.NotNil(t, configurator)

	config, err := configurator.Parse()
	assert.Nil(t, err)
	assert.NotNil(t, config)
}
