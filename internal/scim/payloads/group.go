package payloads

import (
	"io"
)

type CreateScimGroupPayload struct {
	Schemas     []string `json:"schemas"`
	DisplayName string   `json:"displayName"`
}

func GroupPayloadFromJSON(r io.Reader) (*CreateScimGroupPayload, error) {
	var payload CreateScimGroupPayload
	err := decodeJSON(r, &payload)
	if err != nil {
		return nil, err
	}

	return &payload, nil
}
