package scim

import (
	"bytes"
	"encoding/json"
	"net/http"
)

func RenderScimJSON(w http.ResponseWriter, _ *http.Request, status int, v interface{}) {
	buf := &bytes.Buffer{}
	enc := json.NewEncoder(buf)
	enc.SetEscapeHTML(true)

	err := enc.Encode(v)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/scim+json")

	w.WriteHeader(status)
	_, _ = w.Write(buf.Bytes())
}
