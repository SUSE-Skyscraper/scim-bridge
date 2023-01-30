package responses

import (
	"net/http"

	"github.com/go-chi/render"
)

var ErrInternalServerError = &ErrResponse{Schemas: errorSchema, HTTPStatusCode: 500, Details: "Internal server error"}
var errorSchema = []string{"urn:ietf:params:scim:api:messages:2.0:Error"}

type ErrResponse struct {
	Schemas        []string `json:"schemas"`
	Details        string   `json:"details"`
	HTTPStatusCode int      `json:"status"`
	ScimType       string   `json:"scimType,omitempty"`
}

func (e *ErrResponse) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, e.HTTPStatusCode)
	return nil
}

func ErrNotFound(id string) render.Renderer {
	details := "Resource " + id + " not found"

	return &ErrResponse{
		Schemas:        errorSchema,
		HTTPStatusCode: 404,
		Details:        details,
	}
}

func ErrBadValue(err error) render.Renderer {
	return &ErrResponse{
		Schemas:        errorSchema,
		ScimType:       "invalidValue",
		Details:        err.Error(),
		HTTPStatusCode: 400,
	}
}

func ErrBadFilter(err error) render.Renderer {
	return &ErrResponse{
		Schemas:        errorSchema,
		ScimType:       "invalidFilter",
		Details:        err.Error(),
		HTTPStatusCode: 400,
	}
}
