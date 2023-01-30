package responses

import (
	"fmt"
	"github.com/suse-skyscraper/openfga-scim-bridge/bridge"
	"github.com/suse-skyscraper/openfga-scim-bridge/database"
	"net/http"
	"time"
)

type ScimGroupResponse struct {
	Schemas     []string            `json:"schemas,omitempty"`
	ID          string              `json:"id"`
	DisplayName string              `json:"displayName"`
	Members     []map[string]string `json:"members"`
	Meta        map[string]string   `json:"meta"`
}

func (rd *ScimGroupResponse) Render(_ http.ResponseWriter, _ *http.Request) error {
	return nil
}

type ScimListGroupsResponse struct {
	Schemas      []string             `json:"schemas"`
	ItemsPerPage int                  `json:"itemsPerPage"`
	StartIndex   int                  `json:"startIndex"`
	TotalResults int                  `json:"totalResults"`
	Resources    []*ScimGroupResponse `json:"Resources"`
}

func (rd *ScimListGroupsResponse) Render(_ http.ResponseWriter, _ *http.Request) error {
	return nil
}

func NewScimGroupResponse(
	bridge *bridge.Bridge,
	group database.Group,
	members []database.GroupMembership,
) *ScimGroupResponse {
	var memberships []map[string]string
	for _, member := range members {
		memberships = append(memberships, map[string]string{
			"value":   member.UserID.String(),
			"display": member.Username.String,
		})
	}
	return newScimGroupResponse(bridge, group, memberships, true)
}

func newScimGroupResponse(
	bridge *bridge.Bridge,
	group database.Group,
	members []map[string]string,
	singleResponse bool,
) *ScimGroupResponse {
	// the schemas should be added if the response is a single user, not a list
	var schemas []string
	if singleResponse {
		schemas = []string{"urn:ietf:params:scim:schemas:core:2.0:Group"}
	}

	return &ScimGroupResponse{
		Schemas:     schemas,
		ID:          group.ID.String(),
		DisplayName: group.DisplayName,
		Members:     members,
		Meta: map[string]string{
			"resourceType": "Group",
			"created":      group.CreatedAt.Format(time.RFC3339),
			"lastModified": group.UpdatedAt.Format(time.RFC3339),
			"location":     fmt.Sprintf("%s/scim/v2/Groups/%s", bridge.BaseURL, group.ID.String()),
		},
	}
}

type ScimGroupListResponseInput struct {
	TotalResults int
	StartIndex   int
	ItemsPerPage int
}

func NewScimGroupListResponse(
	bridge *bridge.Bridge,
	groups []database.Group,
	input ScimGroupListResponseInput,
) *ScimListGroupsResponse {
	var list []*ScimGroupResponse
	for _, group := range groups {
		list = append(list, newScimGroupResponse(bridge, group, []map[string]string{}, false))
	}

	return &ScimListGroupsResponse{
		Schemas:      []string{"urn:ietf:params:scim:api:messages:2.0:ListResponse"},
		Resources:    list,
		TotalResults: input.TotalResults,
		StartIndex:   input.StartIndex,
		ItemsPerPage: input.ItemsPerPage,
	}
}
