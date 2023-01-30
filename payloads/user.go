package payloads

import (
	"io"
)

type UserEmail struct {
	Value   string `json:"value"`
	Type    string `json:"type"`
	Primary bool   `json:"primary"`
}

type CreateScimUserPayload struct {
	Schemas           []string          `json:"schemas"`
	Username          string            `json:"userName"`
	ID                string            `json:"id"`
	ExternalID        string            `json:"externalId"`
	Name              map[string]string `json:"name"`
	Emails            []UserEmail       `json:"emails"`
	Nickname          string            `json:"nickName"`
	DisplayName       string            `json:"displayName"`
	Active            bool              `json:"active"`
	PreferredLanguage string            `json:"preferredLanguage"`
	Locale            string            `json:"locale"`
	Title             string            `json:"title"`
	UserType          string            `json:"userType"`
	Timezone          string            `json:"timezone"`
}

func CreateUserPayloadFromJSON(r io.Reader) (*CreateScimUserPayload, error) {
	payload := CreateScimUserPayload{
		Active: true,
	}
	err := decodeJSON(r, &payload)
	if err != nil {
		return nil, err
	}

	return &payload, nil
}
