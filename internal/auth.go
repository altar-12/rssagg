package auth

import (
	"errors"
	"net/http"
	"strings"
)

// Extracts the API key from the request header
// header format is `Authorization: ApiKey {insert your api key here}`
func GetAPIKey(h http.Header) (string, error) {
	val := h.Get("Authorization")
	if val == "" {
		return "", errors.New("no authorization header speified")
	}
	vals := strings.Split(val, " ")
	if len(vals) != 2 {
		return "", errors.New("malformed authorization header specified")
	}
	if vals[0] != "ApiKey" {
		return "", errors.New("malformed authorization header specified")
	}
	return vals[1], nil
}
