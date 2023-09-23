package helpers

import (
	"encoding/json"
	"net/http"
)

func RespondWithJSON(w http.ResponseWriter, status int, payload interface{}) error {
	data, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	w.WriteHeader(status)
	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
	return nil
}

func RespondWithError(w http.ResponseWriter, code int, msg string) error {
	return RespondWithJSON(w, code, map[string]string{"error": msg})
}

func DecodeBodyToJson[T any](r *http.Request, data *T) error {
	decoder := json.NewDecoder(r.Body)
	return decoder.Decode(&data)
}
