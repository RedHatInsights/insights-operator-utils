package helpers_test

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/RedHatInsights/insights-operator-utils/tests/helpers"
)

const (
	localhostAddress = "http://localhost"
	port             = 9999
	notJSONString    = ""
)

var (
	testError = fmt.Errorf("test error")
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
