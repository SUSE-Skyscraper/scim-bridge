package application

import (
	"context"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestApp_Start(t *testing.T) {
	ctx := context.Background()
	configDir, err := filepath.Abs("../..")
	assert.Nil(t, err)

	app, err := NewApp(configDir)
	assert.Nil(t, err)

	err = app.Start(ctx)
	assert.Nil(t, err)
}

func TestApp_Start_MissingConfig(t *testing.T) {
	ctx := context.Background()
	configDir, err := filepath.Abs("../../testdata/missing-config")
	assert.Nil(t, err)

	app, err := NewApp(configDir)
	assert.NotNil(t, err)

	err = app.Start(ctx)
	assert.NotNil(t, err)
}

func TestApp_Start_BadDB(t *testing.T) {
	ctx := context.Background()
	configDir, err := filepath.Abs("../..")
	assert.Nil(t, err)

	app, err := NewApp(configDir)
	assert.Nil(t, err)
	app.Config.DB.Host = "badhost"

	err = app.Start(ctx)
	assert.NotNil(t, err)
}

func TestApp_Start_NoDB(t *testing.T) {
	ctx := context.Background()
	configDir, err := filepath.Abs("../..")
	assert.Nil(t, err)

	app, err := NewApp(configDir)
	assert.Nil(t, err)

	app.Config.DB = DBConfig{}

	err = app.Start(ctx)
	assert.NotNil(t, err)
}

func TestApp_Shutdown(t *testing.T) {
	ctx := context.Background()
	configDir, err := filepath.Abs("../..")
	assert.Nil(t, err)

	app, err := NewApp(configDir)
	assert.Nil(t, err)

	err = app.Start(ctx)
	assert.Nil(t, err)

	app.Shutdown(ctx)
}
