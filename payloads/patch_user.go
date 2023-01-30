package payloads

import "io"

type UserPatch struct {
	Active bool `json:"active"`
}

type UserPatchOperation struct {
	Op    string    `json:"op"`
	Path  string    `json:"path"`
	Value UserPatch `json:"value"`
}

type UserPatchPayload struct {
	Schemas    []string              `json:"schemas"`
	Operations []*UserPatchOperation `json:"Operations"`
}

func UserPatchPayloadFromJSON(r io.Reader) (*UserPatchPayload, error) {
	var payload UserPatchPayload
	err := decodeJSON(r, &payload)
	if err != nil {
		return nil, err
	}

	return &payload, nil
}
