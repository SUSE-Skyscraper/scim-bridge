package responses

import (
	"fmt"
	"net/http"
	"time"

	openfga_scim_bridge "github.com/suse-skyscraper/openfga-scim-bridge/v2/bridge"
	"github.com/suse-skyscraper/openfga-scim-bridge/v2/database"
	"github.com/suse-skyscraper/openfga-scim-bridge/v2/payloads"
)

type ScimUserResponse struct {
	Schemas  []string             `json:"schemas,omitempty"`
	UserName string               `json:"userName"`
	ID       string               `json:"id"`
	Name     map[string]string    `json:"name,omitempty"`
	Emails   []payloads.UserEmail `json:"emails,omitempty"`
	Active   bool                 `json:"active"`
	Meta     map[string]string    `json:"meta"`
}

func (rd *ScimUserResponse) Render(_ http.ResponseWriter, _ *http.Request) error {
	return nil
}

type ScimListUsersResponse struct {
	Schemas      []string            `json:"schemas"`
	ItemsPerPage int                 `json:"itemsPerPage"`
	StartIndex   int                 `json:"startIndex"`
	TotalResults int                 `json:"totalResults"`
	Resources    []*ScimUserResponse `json:"Resources"`
}

func (rd *ScimListUsersResponse) Render(_ http.ResponseWriter, _ *http.Request) error {
	return nil
}

func NewScimUserResponse(bridge *openfga_scim_bridge.Bridge, user database.User) *ScimUserResponse {
	return newScimUserResponse(bridge, user, true)
}

func newScimUserResponse(bridge *openfga_scim_bridge.Bridge, user database.User, singleResponse bool) *ScimUserResponse {
	// the schemas should be added if the response is a single user, not a list
	var schemas []string
	if singleResponse {
		schemas = []string{"urn:ietf:params:scim:schemas:core:2.0:User"}
	}

	return &ScimUserResponse{
		Schemas:  schemas,
		ID:       user.ID.String(),
		UserName: user.Username,
		Name:     user.Name,
		Emails:   user.Emails,
		Active:   user.Active,
		Meta: map[string]string{
			"resourceType": "User",
			"created":      user.CreatedAt.Format(time.RFC3339),
			"lastModified": user.UpdatedAt.Format(time.RFC3339),
			"location":     fmt.Sprintf("%s/scim/v2/Users/%s", bridge.BaseURL, user.ID),
		},
	}
}

type ScimUserListResponseInput struct {
	TotalResults int
	StartIndex   int
	ItemsPerPage int
}

func NewScimUserListResponse(
	bridge *openfga_scim_bridge.Bridge,
	users []database.User,
	input ScimUserListResponseInput,
) *ScimListUsersResponse {
	var list []*ScimUserResponse
	for _, user := range users {
		list = append(list, newScimUserResponse(bridge, user, false))
	}

	return &ScimListUsersResponse{
		Schemas:      []string{"urn:ietf:params:scim:api:messages:2.0:ListResponse"},
		Resources:    list,
		TotalResults: input.TotalResults,
		StartIndex:   input.StartIndex,
		ItemsPerPage: input.ItemsPerPage,
	}
}
