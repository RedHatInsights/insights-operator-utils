package helpers_test

// Documentation in literate-programming-style is available at:
// https://redhatinsights.github.io/insights-operator-utils/packages/tests/helpers/helpers_test.html

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/RedHatInsights/insights-results-aggregator-data/testdata"
	"github.com/golang/mock/gomock"
	"github.com/mozillazg/request"
	"github.com/stretchr/testify/assert"
	"github.com/verdverm/frisby"
	"gopkg.in/h2non/gock.v1"

	httputils "github.com/RedHatInsights/insights-operator-utils/http"
	"github.com/RedHatInsights/insights-operator-utils/responses"
	"github.com/RedHatInsights/insights-operator-utils/tests/helpers"
	"github.com/RedHatInsights/insights-operator-utils/tests/mock_io"
	"github.com/RedHatInsights/insights-operator-utils/types"
)

const (
	localhostAddress = "http://localhost"
	port             = 9999
	notJSONString    = "not-json"
	okBody           = `{"status": "ok"}`
	testEndpoint     = "test"
)

var (
	serverAddress = localhostAddress + ":" + fmt.Sprint(port)
	testError     = fmt.Errorf("test error")
	devNull       interface{}
)

func TestFailOnError(t *testing.T) {
	helpers.FailOnError(t, nil)
}

func TestFailOnError_Fatal(t *testing.T) {
	mockT := helpers.NewMockT(t)
	defer mockT.Finish()

	mockT.ExpectFailOnError(testError)

	helpers.FailOnError(mockT, testError)
}

func TestToJSONString(t *testing.T) {
	assert.Equal(t, `{"test":1}`, helpers.ToJSONString(map[string]int{
		"test": 1,
	}))
}

func TestToJSONString_Error(t *testing.T) {
	assert.Panics(t, func() {
		helpers.ToJSONString(make(chan int))
	}, "should panic on unsupported type")
}

func TestToJSONPrettyString(t *testing.T) {
	helpers.AssertStringsAreEqualJSON(t, `{"test": 1, "k": 2}`, helpers.ToJSONPrettyString(map[string]int{
		"test": 1,
		"k":    2,
	}))
}

func TestNewMicroHTTPServer(t *testing.T) {
	server := helpers.NewMicroHTTPServer(":"+fmt.Sprint(port), "")
	_ = server.Initialize()
	server.AddEndpoint("/", func(http.ResponseWriter, *http.Request) {})
}

func TestMustGobSerialize(t *testing.T) {
	objectToSerialize := 1
	bytesResult := helpers.MustGobSerialize(t, objectToSerialize)
	expectedBytes := []byte{0x3, 0x4, 0x0, 0x2}

	assert.Equal(t, expectedBytes, bytesResult)
}

func TestAssertStringsAreEqualJSON(t *testing.T) {
	helpers.AssertStringsAreEqualJSON(t, `{"one": 1, "two": 2}`, `{"two": 2, "one": 1}`)
}

func TestAssertStringsAreEqualJSON_Error(t *testing.T) {
	t.Run("ExpectedIsNotJSON", func(t *testing.T) {
		mockT := helpers.NewMockT(t)
		defer mockT.Finish()

		mockT.ExpectFailOnErrorAnyArgument()
		mockT.Expects.EXPECT().Errorf(gomock.Any(), gomock.Any())

		helpers.AssertStringsAreEqualJSON(mockT, notJSONString, `{"two": 2, "one": 1}`)
	})
	t.Run("GotIsNotJSON", func(t *testing.T) {
		mockT := helpers.NewMockT(t)
		defer mockT.Finish()

		mockT.ExpectFailOnErrorAnyArgument()
		mockT.Expects.EXPECT().Errorf(gomock.Any(), gomock.Any())

		helpers.AssertStringsAreEqualJSON(mockT, `{"one": 1, "two": 2}`, notJSONString)
	})
}

func TestJSONUnmarshalStrict(t *testing.T) {
	jsonBytes := []byte(`{"one": 1, "two": 2, "three": 3}`)
	var resultObj map[string]int

	err := helpers.JSONUnmarshalStrict(jsonBytes, &resultObj)
	helpers.FailOnError(t, err)

	assert.Equal(t, map[string]int{"one": 1, "two": 2, "three": 3}, resultObj)
}

func TestIsStringJSON(t *testing.T) {
	t.Run("true", func(t *testing.T) {
		assert.True(t, helpers.IsStringJSON(`{"one": 1, "two": 2}`))
	})
	t.Run("false", func(t *testing.T) {
		assert.False(t, helpers.IsStringJSON(`{"one": 1"two": 2}`))
	})
}

func TestRunTestWithTimeout(t *testing.T) {
	helpers.RunTestWithTimeout(t, func(t testing.TB) {}, time.Second)
}

func TestRunTestWithTimeout_Error(t *testing.T) {
	mockT := helpers.NewMockT(t)
	defer mockT.Finish()

	mockT.Expects.EXPECT().Fatal("test ran out of time")

	helpers.RunTestWithTimeout(mockT, func(t testing.TB) {
		time.Sleep(time.Hour)
	}, time.Microsecond)
}

func TestAssertAPIRequest(t *testing.T) {
	const (
		apiPrefix           = "/api/v1/"
		endpoint            = "test/{param}"
		expectedURLParam    = uint64(55)
		expectedRequestBody = `{"test": "json"}`
	)

	token := helpers.MakeXRHTokenString(t, &types.Token{
		Identity: types.Identity{
			AccountNumber: testdata.UserID,
			Internal: types.Internal{
				OrgID: testdata.OrgID,
			},
		},
	})

	testServer := helpers.NewMicroHTTPServer(":8080", apiPrefix)
	testServer.AddEndpoint(endpoint, func(writer http.ResponseWriter, request *http.Request) {
		paramValue, err := httputils.GetRouterPositiveIntParam(request, "param")
		helpers.FailOnError(t, err)

		assert.Equal(t, expectedURLParam, paramValue)

		helpers.CheckResponseBodyJSON(t, expectedRequestBody, request.Body)

		err = responses.SendOK(writer, map[string]interface{}{"param": paramValue})
		helpers.FailOnError(t, err)
	})

	helpers.AssertAPIRequest(t, testServer, apiPrefix, &helpers.APIRequest{
		Method:             http.MethodGet,
		Endpoint:           endpoint,
		EndpointArgs:       []interface{}{expectedURLParam},
		Body:               expectedRequestBody,
		UserID:             testdata.UserID,
		OrgID:              testdata.OrgID,
		XRHIdentity:        token,
		AuthorizationToken: token,
		ExtraHeaders: map[string][]string{
			"User-Agent": {"test"},
		},
	}, &helpers.APIResponse{
		StatusCode: http.StatusOK,
		Body:       fmt.Sprintf(`{"param": %v}`, expectedURLParam),
		Headers: map[string]string{
			// should be empty, which means there's no such header
			"non-existing-header": "",
		},
		BodyChecker: func(t testing.TB, expected, got []byte) {
			helpers.AssertStringsAreEqualJSON(t, string(expected), string(got))
		},
	})

	t.Run("NoBodyChecker", func(t *testing.T) {
		helpers.AssertAPIRequest(t, testServer, apiPrefix, &helpers.APIRequest{
			Method:       http.MethodGet,
			Endpoint:     endpoint,
			EndpointArgs: []interface{}{expectedURLParam},
			Body:         expectedRequestBody,
		}, &helpers.APIResponse{
			StatusCode: http.StatusOK,
			Body:       fmt.Sprintf(`{"param": %v}`, expectedURLParam),
		})
	})
}

func TestAssertReportResponsesEqual(t *testing.T) {
	testAssertResponsesEqual(
		t,
		helpers.AssertReportResponsesEqual,
		testdata.Report3RulesExpectedResponse,
		5,
		4,
	)
}

func TestAssertRuleResponsesEqual(t *testing.T) {
	correctResponse := struct {
		Status string             `json:"status"`
		Report types.RuleOnReport `json:"report"`
	}{
		Status: "ok",
	}

	testAssertResponsesEqual(
		t,
		helpers.AssertRuleResponsesEqual,
		helpers.ToJSONString(correctResponse),
		2,
		1,
	)
}

func testAssertResponsesEqual(
	t *testing.T,
	function func(testing.TB, []byte, []byte),
	correctResponse string,
	numberOfErrorsOnBadExpectedValue int,
	numberOfErrorsOnBadGotValue int,
) {
	t.Run("OK", func(t *testing.T) {
		function(
			t,
			[]byte(correctResponse),
			[]byte(correctResponse),
		)
	})

	mockT := helpers.NewMockT(t)
	defer mockT.Finish()

	var devNull interface{}
	badJSONError := json.Unmarshal([]byte(notJSONString), &devNull)

	expectAssertReportResponsesEqualWrongArg := func(n int) {
		mockT.ExpectFailOnError(badJSONError)
		// each AssertReportResponsesEqual raises a lot of errors
		for i := 0; i < n; i++ {
			mockT.Expects.EXPECT().Errorf(gomock.Any(), gomock.Any())
		}
	}

	t.Run("BadExpectedValue", func(t *testing.T) {
		expectAssertReportResponsesEqualWrongArg(numberOfErrorsOnBadExpectedValue)
		function(
			mockT,
			[]byte(notJSONString),
			[]byte(correctResponse),
		)
	})

	t.Run("BadGotValue", func(t *testing.T) {
		expectAssertReportResponsesEqualWrongArg(numberOfErrorsOnBadGotValue)
		function(
			mockT,
			[]byte(correctResponse),
			[]byte("not-json"),
		)
	})
}

func TestGockExpectAPIRequest(t *testing.T) {
	defer helpers.CleanAfterGock(t)

	helpers.GockExpectAPIRequest(t, serverAddress, &helpers.APIRequest{
		Method:   http.MethodPost,
		Endpoint: testEndpoint,
		Body:     okBody,
		ExtraHeaders: map[string][]string{
			"Content-Type": {"application/json"},
		},
	}, &helpers.APIResponse{
		StatusCode: http.StatusOK,
		Body:       okBody,
	})

	resp, err := http.Post(serverAddress+"/"+testEndpoint, "application/json", strings.NewReader(okBody))
	helpers.FailOnError(t, err)

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	helpers.FailOnError(t, err)

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, okBody, string(bodyBytes))
}

func TestNewGockAPIEndpointMatcher(t *testing.T) {
	matcher := helpers.NewGockAPIEndpointMatcher(testEndpoint)

	request, err := http.NewRequest(http.MethodGet, testEndpoint, nil)
	helpers.FailOnError(t, err)

	result, err := matcher(request, nil)
	helpers.FailOnError(t, err)

	assert.True(t, result)
}

func TestCleanAfterGockError(t *testing.T) {
	mockT := helpers.NewMockT(t)
	defer mockT.Finish()

	// remove is crucial
	gock.Remove(gock.New(serverAddress).Mock)
	defer helpers.CleanAfterGock(mockT)

	_, err := http.Post(serverAddress, "application/json", strings.NewReader("{}"))
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "gock: cannot match any request")

	mockT.Expects.EXPECT().Error(gomock.Any())
	for i := 0; i < 4; i++ {
		mockT.Expects.EXPECT().Errorf(gomock.Any(), gomock.Any())
	}
	mockT.Expects.EXPECT().Fatalf(gomock.Any(), gomock.Any())
}

func TestToBytes(t *testing.T) {
	t.Run("OK", func(t *testing.T) {
		assert.Equal(t, []byte(nil), helpers.ToBytes(t, nil))
		const testStr = "1"
		assert.Equal(t, []byte(testStr), helpers.ToBytes(t, []byte(testStr)))
		assert.Equal(t, []byte(testStr), helpers.ToBytes(t, testStr))
		assert.Equal(t, []byte(testStr), helpers.ToBytes(t, strings.NewReader(testStr)))
	})

	t.Run("Error", func(t *testing.T) {
		mockT := helpers.NewMockT(t)
		defer mockT.Finish()

		mockT.Expects.EXPECT().Fatalf(gomock.Any(), gomock.Any())

		_ = helpers.ToBytes(mockT, []string{})
	})
}

func TestUnmarshalResponseBodyToJSON(t *testing.T) {
	var result map[string]int

	err := helpers.UnmarshalResponseBodyToJSON(
		ioutil.NopCloser(strings.NewReader(`{"test": 1}`)),
		&result,
	)
	helpers.FailOnError(t, err)

	assert.Len(t, result, 1)
	assert.Equal(t, 1, result["test"])
}

func TestUnmarshalResponseBodyToJSON_UnmarshalError(t *testing.T) {
	err := helpers.UnmarshalResponseBodyToJSON(
		ioutil.NopCloser(strings.NewReader(notJSONString)),
		&devNull,
	)
	assert.EqualError(t, err, "invalid character 'o' in literal null (expecting 'u')")
}

func TestUnmarshalResponseBodyToJSON_ReaderError(t *testing.T) {
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	mockReadCloser := mock_io.NewMockReadCloser(mockController)

	t.Run("ReadError", func(t *testing.T) {
		mockReadCloser.EXPECT().Read(gomock.Any()).Return(0, testError)

		err := helpers.UnmarshalResponseBodyToJSON(
			mockReadCloser,
			&devNull,
		)
		assert.EqualError(t, err, testError.Error())
	})

	t.Run("CloseError", func(t *testing.T) {
		mockReadCloser.EXPECT().Read(gomock.Any()).Return(0, io.EOF)
		mockReadCloser.EXPECT().Close().Return(testError)

		err := helpers.UnmarshalResponseBodyToJSON(
			mockReadCloser,
			&devNull,
		)
		assert.EqualError(t, err, testError.Error())
	})
}

func TestFrisbyExpectItemInArray(t *testing.T) {
	testFrisbyExpectItemInArray(
		t,
		`{"organizations": [1, 2, 3, 55], "status": "ok"}`,
		true,
		"",
	)
	testFrisbyExpectItemInArray(
		t,
		notJSONString,
		false,
		"invalid character 'o' in literal null (expecting 'u')",
	)
	testFrisbyExpectItemInArray(
		t,
		`{"orgs": [1, 2, 3, 55], "status": "ok"}`,
		false,
		`field organizations does not exist in response {"orgs":[1,2,3,55],"status":"ok"}`,
	)
	testFrisbyExpectItemInArray(
		t,
		`{"organizations": 1, "status": "ok"}`,
		false,
		`field organizations is not an array in response {"organizations":1,"status":"ok"}`,
	)
	testFrisbyExpectItemInArray(
		t,
		`{"organizations": [1, 2, 3], "status": "ok"}`,
		false,
		`Item 55 was not found in array [1 2 3] in response {"organizations":[1,2,3],"status":"ok"}`,
	)
}

func testFrisbyExpectItemInArray(t *testing.T, responseBody string, expectedFound bool, expectedErrStr string) {
	checker := helpers.FrisbyExpectItemInArray("organizations", 55)

	f := frisby.Create("")
	f.Resp = &request.Response{
		Response: &http.Response{
			Status:        "ok",
			StatusCode:    http.StatusOK,
			Body:          ioutil.NopCloser(strings.NewReader(responseBody)),
			ContentLength: int64(len(responseBody)),
		},
	}

	found, errStr := checker(f)
	assert.Equal(t, expectedFound, found)
	assert.Equal(t, expectedErrStr, errStr)
}
