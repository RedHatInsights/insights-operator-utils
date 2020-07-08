package httputils

import (
	"fmt"
	"net/url"
	"regexp"
	"strings"
)

// MakeURLToEndpoint creates URL to endpoint, use constants from file endpoints.go
func MakeURLToEndpoint(apiPrefix, endpoint string, args ...interface{}) string {
	endpoint = ReplaceParamsInEndpointAndTrimLeftSlash(endpoint, "%v")

	apiPrefix = strings.TrimRight(apiPrefix, "/")
	nonParsedURL := apiPrefix + "/" + fmt.Sprintf(endpoint, args...)
	resultingURL, err := url.Parse(nonParsedURL)

	if err != nil {
		return nonParsedURL
	}

	return resultingURL.String()
}

// ReplaceParamsInEndpointAndTrimLeftSlash replaces params in endpoint and trims left slash
func ReplaceParamsInEndpointAndTrimLeftSlash(endpoint, replacer string) string {
	re := regexp.MustCompile(`\{[a-zA-Z_0-9]+\}`)

	endpoint = re.ReplaceAllString(endpoint, replacer)
	endpoint = strings.TrimLeft(endpoint, "/")

	return endpoint
}

// MakeURLToEndpointMapString creates URL to endpoint using arguments in map in string format, use constants from file endpoints.go
func MakeURLToEndpointMapString(apiPrefix, endpoint string, args map[string]string) string {
	newArgs := make(map[string]interface{})

	for key, val := range args {
		newArgs[key] = val
	}

	return MakeURLToEndpointMap(apiPrefix, endpoint, newArgs)
}

// MakeURLToEndpointMap creates URL to endpoint using arguments in map, use constants from file endpoints.go
func MakeURLToEndpointMap(apiPrefix, endpoint string, args map[string]interface{}) string {
	endpoint = strings.TrimLeft(endpoint, "/")
	for key, val := range args {
		endpoint = strings.ReplaceAll(endpoint, fmt.Sprintf("{%v}", key), fmt.Sprint(val))
	}

	apiPrefix = strings.TrimRight(apiPrefix, "/")

	return apiPrefix + "/" + endpoint
}
