package middleware

import (
	"fmt"
	"net/http"

	"github.com/go-chi/render"
	"github.com/suse-skyscraper/openfga-scim-bridge/internal/application"
	"github.com/suse-skyscraper/openfga-scim-bridge/internal/auth/apikeys"
	"github.com/suse-skyscraper/openfga-scim-bridge/internal/scim/responses"
)

func BearerAuthorizationHandler(app *application.App) func(next http.Handler) http.Handler {
	verifier := apikeys.NewVerifier(app)

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authorizationHeader := r.Header.Get("Authorization")
			match, err := verifier.VerifyScim(r.Context(), authorizationHeader)
			if err != nil {
				_ = render.Render(w, r, responses.ErrInternalServerError)
				return
			} else if !match {
				w.WriteHeader(http.StatusUnauthorized)
				_, _ = fmt.Fprintf(w, "Not Authorized")
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
