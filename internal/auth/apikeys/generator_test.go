package apikeys

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerator_Generate(t *testing.T) {
	g := Generator{
		Memory:      64 * 1024,
		Time:        1,
		Parallelism: 4,
	}
	apiKey, hash, err := g.Generate()
	assert.Nil(t, err)
	assert.NotEqualf(t, "", apiKey, "apiKey should not be empty")
	assert.NotEqualf(t, "", hash, "hash should not be empty")
}
