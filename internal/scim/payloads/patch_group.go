package payloads

import (
	"errors"
	"io"
	"regexp"

	"github.com/google/uuid"
)

type GroupPatchOperation struct {
	Op    string      `json:"op"`
	Path  string      `json:"path"`
	Value interface{} `json:"value"`
}

type GroupPatchPayload struct {
	Schemas    []string               `json:"schemas"`
	Operations []*GroupPatchOperation `json:"Operations"`
}

func GroupPatchPayloadFromJSON(r io.Reader) (*GroupPatchPayload, error) {
	var payload GroupPatchPayload
	err := decodeJSON(r, &payload)
	if err != nil {
		return nil, err
	}

	return &payload, nil
}

type MemberPatch struct {
	Display string
	Value   uuid.UUID
}

type GroupPatch struct {
	DisplayName string
}

func (o *GroupPatchOperation) ParseIDFromPath() (uuid.UUID, error) {
	r := regexp.MustCompile(`^members\[value eq "(\S+)"]$`)
	match := r.FindStringSubmatch(o.Path)
	if match == nil {
		return uuid.UUID{}, errors.New("invalid path")
	}

	idString := match[1]
	id, err := uuid.Parse(idString)
	if err != nil {
		return uuid.UUID{}, errors.New("invalid value type")
	}

	return id, nil
}

func (o *GroupPatchOperation) GetPatch() (GroupPatch, error) {
	values, ok := o.Value.(map[string]interface{})
	if !ok {
		return GroupPatch{}, errors.New("invalid value type")
	}

	displayName, ok := values["displayName"].(string)
	if !ok {
		return GroupPatch{}, errors.New("invalid value type")
	}

	patch := GroupPatch{
		DisplayName: displayName,
	}

	return patch, nil
}

func (o *GroupPatchOperation) GetAddMembersPatch() ([]uuid.UUID, error) {
	members, ok := o.Value.([]interface{})
	if !ok {
		return []uuid.UUID{}, errors.New("invalid value type")
	}

	var patches []uuid.UUID

	for _, memberInterface := range members {
		member, ok := memberInterface.(map[string]interface{})
		if !ok {
			return []uuid.UUID{}, errors.New("invalid value type")
		}

		value, ok := member["value"].(string)
		if !ok {
			return []uuid.UUID{}, errors.New("invalid value type")
		}

		id, err := uuid.Parse(value)
		if err != nil {
			return []uuid.UUID{}, errors.New("invalid value type")
		}

		patches = append(patches, id)
	}

	return patches, nil
}
