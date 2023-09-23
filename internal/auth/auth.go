package auth

import (
	"errors"
	"net/http"
	"strings"
)

func getAuthToken(r *http.Request, strip string) (string, error) {
	authHeader := strings.Replace(r.Header.Get("Authorization"), strip, "", 1)
	if len(authHeader) == 0 {
		return "", errors.New("No authorization header recieved")
	}
	return authHeader, nil
}

// Example
// Authorization: Bearer <access_token>
func GetAuthBearer(r *http.Request) (string, error) {
	return getAuthToken(r, "Bearer ")
}

// Example
// Authorization: ApiKey <api_key>
func GetAuthApiKey(r *http.Request) (string, error) {
	return getAuthToken(r, "ApiKey ")
}
