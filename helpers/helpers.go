package helpers

import (
	"encoding/json"
	"fmt"
	"io"
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

func GetXML(url string) ([]byte, error) {
	fmt.Println("Fetching XML from: ", url)
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Fetching ERROR:", err)
		return []byte{}, fmt.Errorf("GET error: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return []byte{}, fmt.Errorf("Status error: %v", resp.StatusCode)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return []byte{}, fmt.Errorf("Read body: %v", err)
	}

	return data, nil
}
