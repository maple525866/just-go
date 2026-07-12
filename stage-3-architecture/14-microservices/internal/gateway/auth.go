package gateway

import (
	"crypto/subtle"
	"strings"
)

func authorizedBearer(header, expectedToken string) bool {
	expectedToken = strings.TrimSpace(expectedToken)
	if expectedToken == "" {
		return false
	}
	expected := "Bearer " + expectedToken
	return subtle.ConstantTimeCompare([]byte(header), []byte(expected)) == 1
}
