package helpers

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"regexp"
	"strings"
	"testing"

	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"
	gock "gopkg.in/h2non/gock.v1"

	httputils "github.com/RedHatInsights/insights-operator-utils/http"
	"github.com/RedHatInsights/insights-operator-utils/types"
)

// ServerInitializer is interface which is implemented for any server having Initialize method
type ServerInitializer interface {
	Initialize() http.Handler
}

// BodyChecker represents body checker type for api response
type BodyChecker func(t testing.TB, expected, got []byte)

// APIRequest is a request to api to use in AssertAPIRequest
//
// (required) Method is an http method
// (required) Endpoint is an endpoint without api prefix
// EndpointArgs are the arguments to pass to endpoint template (leave empty if endpoint is not a template)
// Body is a request body which can be a string or []byte (leave empty to not send)
// UserID is a user id for methods requiring user id (leave empty to not use it)
// OrgID is an org id for methods requiring it to be in token (leave empty to not use it)
// XRHIdentity is an authentication token (leave empty to not use it)
// AuthorizationToken is an authentication token (leave empty to not use it)
type APIRequest struct {
	Method             string
	Endpoint           string
	EndpointArgs       []interface{}
	Body               interface{}
	UserID             types.UserID
	OrgID              types.OrgID
	XRHIdentity        string
	AuthorizationToken string
	ExtraHeaders       http.Header
}

// APIResponse is an expected api response to use in AssertAPIRequest
//
// StatusCode is an expected http status code (leave empty to not check for status code)
// Body is an expected body which can be a string or []byte(leave empty to not check for body)
// BodyChecker is a custom body checker function (leave empty to use default one - CheckResponseBodyJSON)
type APIResponse struct {
	StatusCode  int
	Body        interface{}
	BodyChecker BodyChecker
	Headers     map[string]string
}

// AssertAPIRequest sends sends api request and checks api response (see docs for APIRequest and APIResponse)
// to the provided testServer using the provided APIPrefix
func AssertAPIRequest(
	t testing.TB,
	testServer ServerInitializer,
	APIPrefix string,
	request *APIRequest,
	expectedResponse *APIResponse,

) {
	url := httputils.MakeURLToEndpoint(APIPrefix, request.Endpoint, request.EndpointArgs...)

	req := makeRequest(t, request, url)

	response := ExecuteRequest(testServer, req).Result()

	if len(expectedResponse.Headers) != 0 {
		checkResponseHeaders(t, expectedResponse.Headers, response.Header)
	}
	if expectedResponse.StatusCode != 0 {
		assert.Equal(t, expectedResponse.StatusCode, response.StatusCode, "Expected different status code")
	}

	if expectedResponse.Body != nil {
		assertBody(t, expectedResponse.Body, response.Body, expectedResponse.BodyChecker)
	}
}

func toBytes(t testing.TB, obj interface{}) []byte {
	if obj == nil {
		return nil
	}

	switch v := obj.(type) {
	case []byte:
		return v
	case string:
		return []byte(v)
	case io.Reader:
		res, err := ioutil.ReadAll(v)
		FailOnError(t, err)
		return res
	default:
		t.Fatalf(
			`type "%T" of API(Request|Response).Body is not supported, please see the documentation. Value is "%+v"`,
			obj, obj,
		)
	}

	return nil
}

func assertBody(t testing.TB, expectedBody interface{}, body interface{}, bodyChecker BodyChecker) {
	expectedBodyBytes := toBytes(t, expectedBody)
	bodyBytes := toBytes(t, body)

	if bodyChecker != nil {
		bodyChecker(t, expectedBodyBytes, bodyBytes)
	} else {
		AssertStringsAreEqualJSON(t, string(expectedBodyBytes), string(bodyBytes))
	}
}

// ExecuteRequest executes http request on a testServer
func ExecuteRequest(testServer ServerInitializer, req *http.Request) *httptest.ResponseRecorder {
	router := testServer.Initialize()

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	return rr
}

func makeRequest(t testing.TB, request *APIRequest, url string) *http.Request {
	bodyBytes := toBytes(t, request.Body)

	req := httptest.NewRequest(request.Method, url, bytes.NewReader(bodyBytes))

	// authorize user
	if request.UserID != types.UserID("") || request.OrgID != types.OrgID(0) {
		identity := types.Identity{
			AccountNumber: request.UserID,
			Internal: types.Internal{
				OrgID: request.OrgID,
			},
		}
		req = req.WithContext(context.WithValue(req.Context(), types.ContextKeyUser, identity))
	}

	if len(request.XRHIdentity) != 0 {
		req.Header.Set("x-rh-identity", request.XRHIdentity)
	}

	if len(request.AuthorizationToken) != 0 {
		req.Header.Set("Authorization", request.AuthorizationToken)
	}

	for headerKey, headerValues := range request.ExtraHeaders {
		for _, headerValue := range headerValues {
			req.Header.Add(headerKey, headerValue)
		}
	}

	return req
}

// CheckResponseBodyJSON checks if body is the same json as in expected
// (ignores whitespaces, newlines, etc)
// also validates both expected and body to be a valid json
func CheckResponseBodyJSON(t testing.TB, expectedJSON string, body io.ReadCloser) {
	result, err := ioutil.ReadAll(body)
	FailOnError(t, err)

	AssertStringsAreEqualJSON(t, expectedJSON, string(result))
}

// checkResponseHeaders checks if headers are the same as in expected
func checkResponseHeaders(t testing.TB, expectedHeaders map[string]string, actualHeaders http.Header) {
	for key, value := range expectedHeaders {
		assert.Equal(t, value, actualHeaders.Get(key), "Expected different headers")
	}
}

// AssertReportResponsesEqual checks if reports in answer are the same
func AssertReportResponsesEqual(t testing.TB, expected, got []byte) {
	var expectedResponse, gotResponse struct {
		Status string               `json:"status"`
		Report types.ReportResponse `json:"report"`
	}

	err := JSONUnmarshalStrict(expected, &expectedResponse)
	if err != nil {
		log.Error().Msg("Error unmarshalling expected value")
	}

	FailOnError(t, err)
	err = JSONUnmarshalStrict(got, &gotResponse)
	if err != nil {
		log.Error().Msg("Error unmarshalling got value")
	}
	FailOnError(t, err)

	assert.NotEmpty(
		t,
		expectedResponse.Status,
		"status is empty(probably json is completely wrong and unmarshal didn't do anything useful)",
	)
	assert.Equal(t, expectedResponse.Status, gotResponse.Status)
	assert.Equal(t, expectedResponse.Report.Meta, gotResponse.Report.Meta)
	// ignore the order
	assert.Equal(
		t,
		len(expectedResponse.Report.Report),
		len(gotResponse.Report.Report),
		"length of reports should be equal",
	)
	assert.ElementsMatch(t, expectedResponse.Report.Report, gotResponse.Report.Report)
}

// NewGockAPIEndpointMatcher returns new matcher for github.com/h2non/gock to match endpoint with any args
func NewGockAPIEndpointMatcher(endpoint string) func(req *http.Request, _ *gock.Request) (bool, error) {
	endpoint = httputils.ReplaceParamsInEndpointAndTrimLeftSlash(endpoint, ".*")
	re := regexp.MustCompile("^" + endpoint + `(\?.*)?$`)

	return func(req *http.Request, _ *gock.Request) (bool, error) {
		uri := req.URL.RequestURI()
		uri = strings.TrimLeft(uri, "/")
		return re.MatchString(uri), nil
	}
}

// NewGockRequestMatcher returns a new matcher for github.com/h2non/gock to match requests
// with provided method, url and body(the same types as body in APIRequest(see the docs))
func NewGockRequestMatcher(
	t *testing.T, method string, url string, body interface{},
) func(*http.Request, *gock.Request) (bool, error) {
	return func(httpReq *http.Request, gockReq *gock.Request) (bool, error) {
		assert.Equal(t, method, httpReq.Method)
		assert.Equal(t, url, httpReq.URL.String())
		if body != nil {
			assertBody(t, body, httpReq.Body, nil)
		}

		return true, nil
	}
}

// GockExpectAPIRequest makes gock expect the request with the baseURL and sends back the response
func GockExpectAPIRequest(t *testing.T, baseURL string, request *APIRequest, response *APIResponse) {
	bodyBytes := toBytes(t, response.Body)

	headers := map[string]string{}

	for key, values := range request.ExtraHeaders {
		for _, value := range values {
			headers[key] = value
		}
	}
	gock.New(baseURL).
		AddMatcher(NewGockRequestMatcher(
			t,
			request.Method,
			httputils.MakeURLToEndpoint(baseURL, request.Endpoint, request.EndpointArgs...),
			request.Body,
		)).
		MatchHeaders(headers).
		Reply(response.StatusCode).
		SetHeaders(response.Headers).
		Body(bytes.NewBuffer(bodyBytes))
}

// CleanAfterGock cleans after gock library and prints all unmatched requests
func CleanAfterGock(t testing.TB) {
	defer gock.Off()
	defer func() {
		hasUnmatchedRequest := false

		for _, request := range gock.GetUnmatchedRequests() {
			hasUnmatchedRequest = false
			fmt.Println("Not expected request: ")

			fmt.Printf("\tMethod: `%+v`\n", request.Method)
			fmt.Printf("\tURL: `%+v`\n", request.URL)
			fmt.Printf("\tHeader: `%+v`\n", ToJSONString(request.Header))

			if request.Body != nil {
				bodyBytes, err := ioutil.ReadAll(request.Body)
				FailOnError(t, err)

				fmt.Printf("\tBody: `%+v`\n", string(bodyBytes))
			}
		}

		if hasUnmatchedRequest {
			t.Fatalf("there were some unexpected requests")
		}
	}()
}
