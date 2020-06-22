package httputils

import (
	"fmt"
	"regexp"
	"strings"
)

// MakeURLToEndpoint creates URL to endpoint, use constants from file endpoints.go
func MakeURLToEndpoint(apiPrefix, endpoint string, args ...interface{}) string {
	endpoint = ReplaceParamsInEndpointAndTrimLeftSlash(endpoint, "%v")

	apiPrefix = strings.TrimRight(apiPrefix, "/")

	return apiPrefix + "/" + fmt.Sprintf(endpoint, args...)
}

// ReplaceParamsInEndpointAndTrimLeftSlash replaces params in endpoint and trims left slash
func ReplaceParamsInEndpointAndTrimLeftSlash(endpoint, replacer string) string {
	re := regexp.MustCompile(`\{[a-zA-Z_0-9]+\}`)

	endpoint = re.ReplaceAllString(endpoint, replacer)
	endpoint = strings.TrimLeft(endpoint, "/")

	return endpoint
}
