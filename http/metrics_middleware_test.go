package httputils_test

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/gorilla/mux"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"

	httputils "github.com/RedHatInsights/insights-operator-utils/http"
	"github.com/RedHatInsights/insights-operator-utils/responses"
	"github.com/RedHatInsights/insights-operator-utils/tests/helpers"
)

const (
	localhostAddress = "http://localhost"
	port             = 9999
)

func TestLogRequest(t *testing.T) {
	buf := new(bytes.Buffer)
	log.Logger = zerolog.New(buf).With().Timestamp().Logger()

	server := createTestServer(t, []Endpoint{
		{
			Path: "/",
			Func: func(writer http.ResponseWriter, request *http.Request) {
				err := responses.Send(http.StatusOK, writer, responses.BuildOkResponse())
				helpers.FailOnError(t, err)
			},
			Methods: []string{http.MethodGet},
		},
	})

	resp, err := http.Get(localhostAddress + ":" + fmt.Sprint(port) + "/")
	helpers.FailOnError(t, err)

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	err = server.Shutdown(context.TODO())
	helpers.FailOnError(t, err)

	assert.Contains(t, buf.String(), "Request received - URI: /, Method: GET")
}

type Endpoint struct {
	Path    string
	Func    func(http.ResponseWriter, *http.Request)
	Methods []string
}

func createTestServer(t testing.TB, endpoints []Endpoint) *http.Server {
	router := mux.NewRouter().StrictSlash(true)
	router.Use(httputils.LogRequest)

	for _, endpoint := range endpoints {
		router.HandleFunc(endpoint.Path, endpoint.Func).Methods(endpoint.Methods...)
	}

	server := &http.Server{Addr: ":" + fmt.Sprint(port), Handler: router}

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			helpers.FailOnError(t, err)
		}
	}()

	return server
}
