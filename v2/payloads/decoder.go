package payloads

import (
	"encoding/json"
	"io"
)

func decodeJSON(r io.Reader, v interface{}) error {
	defer func(dst io.Writer, src io.Reader) {
		_, _ = io.Copy(dst, src)
	}(io.Discard, r)

	return json.NewDecoder(r).Decode(v)
}
