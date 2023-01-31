package middleware

import (
	"fmt"
	"github.com/suse-skyscraper/openfga-scim-bridge/example/internal/apikeys"
	"github.com/suse-skyscraper/openfga-scim-bridge/example/internal/application"
	"net/http"
)

func BearerAuthorizationHandler(app *application.App) func(next http.Handler) http.Handler {
	verifier := apikeys.NewVerifier(app)

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authorizationHeader := r.Header.Get("Authorization")
			match, err := verifier.VerifyScim(r.Context(), authorizationHeader)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				_, _ = fmt.Fprintf(w, "get_text_map_propagator")
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
