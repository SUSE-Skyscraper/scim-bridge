package payloads

import (
	"io"

	"github.com/jackc/pgtype"
)

type ScimEmail struct {
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
	Emails            []ScimEmail       `json:"emails"`
	Nickname          string            `json:"nickName"`
	DisplayName       string            `json:"displayName"`
	Active            bool              `json:"active"`
	PreferredLanguage string            `json:"preferredLanguage"`
	Locale            string            `json:"locale"`
	Title             string            `json:"title"`
	UserType          string            `json:"userType"`
	Timezone          string            `json:"timezone"`
	jsonEmails        pgtype.JSONB
	jsonName          pgtype.JSONB
}

func UserPayloadFromJSON(r io.Reader) (*CreateScimUserPayload, error) {
	payload := CreateScimUserPayload{
		Active: true,
	}
	err := decodeJSON(r, &payload)
	if err != nil {
		return nil, err
	}

	if payload.Name != nil {
		name := pgtype.JSONB{}
		err := name.Set(payload.Name)
		if err != nil {
			return nil, err
		}
		payload.jsonName = name
	} else {
		payload.jsonName = pgtype.JSONB{Bytes: nil, Status: pgtype.Null}
	}

	if payload.Emails != nil {
		emails := pgtype.JSONB{}
		err := emails.Set(payload.Emails)
		if err != nil {
			return nil, err
		}
		payload.jsonEmails = emails
	} else {
		payload.jsonEmails = pgtype.JSONB{Bytes: nil, Status: pgtype.Null}
	}

	return &payload, nil
}

func (u *CreateScimUserPayload) GetJSONName() pgtype.JSONB {
	return u.jsonName
}

func (u *CreateScimUserPayload) GetJSONEmails() pgtype.JSONB {
	return u.jsonEmails
}
