package api

import (
	"encoding/json"
	"net/http"
)

// writeJSON writes v as the JSON response body, mirroring FastAPI's default
// response serialization (200 OK, application/json).
func writeJSON(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

// decodeJSONBody decodes the request body into a generic map, mirroring
// `await request.json()` in the Python handlers.
func decodeJSONBody(r *http.Request) (map[string]interface{}, error) {
	var data map[string]interface{}
	if r.Body == nil {
		return map[string]interface{}{}, nil
	}
	defer r.Body.Close()
	dec := json.NewDecoder(r.Body)
	if err := dec.Decode(&data); err != nil {
		return nil, err
	}
	if data == nil {
		data = map[string]interface{}{}
	}
	return data, nil
}
