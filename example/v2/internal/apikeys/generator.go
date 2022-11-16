package apikeys

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"

	"github.com/suse-skyscraper/openfga-scim-bridge/example/v2/internal/application"
	"golang.org/x/crypto/argon2"
)

const argon2KeyLength = 32

type Generator struct {
	Memory      uint32
	Time        uint32
	Parallelism uint8
}

func NewGenerator(app *application.App) Generator {
	return Generator{
		Memory:      app.Config.Argon2Config.MemoryCost,
		Time:        app.Config.Argon2Config.TimeCost,
		Parallelism: app.Config.Argon2Config.Parallelism,
	}
}

func (g *Generator) Generate() (string, string, error) {
	apiKeyBytes, err := generateRandomBytes(32)
	if err != nil {
		return "", "", err
	}
	apiKey := base64.RawURLEncoding.EncodeToString(apiKeyBytes)

	saltBytes, err := generateRandomBytes(16)
	if err != nil {
		return "", "", err
	}

	hash := argon2.IDKey([]byte(apiKey), saltBytes, g.Time, g.Memory, g.Parallelism, argon2KeyLength)
	b64Salt := base64.RawStdEncoding.EncodeToString(saltBytes)
	b64Hash := base64.RawStdEncoding.EncodeToString(hash)

	encodedHash := fmt.Sprintf("$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s", argon2.Version, g.Memory, g.Time,
		g.Parallelism, b64Salt, b64Hash)

	return encodedHash, apiKey, nil
}

func generateRandomBytes(n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}

	return b, nil
}
