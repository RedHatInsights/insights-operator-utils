package helpers

import (
	"net/http"

	"github.com/gorilla/mux"
)

// MicroHTTPServer in an implementation of ServerInitializer interface
// This small implementation could help implementing tests without using
// a real HTTP server implementation
type MicroHTTPServer struct {
	Serv      *http.Server
	Router    *mux.Router
	APIPrefix string
}

func NewMicroHTTPServer(address string, apiPrefix string) *MicroHTTPServer {
	router := mux.NewRouter().StrictSlash(true)
	server := &http.Server{Addr: address, Handler: router}
	return &MicroHTTPServer{
		APIPrefix: apiPrefix,
		Router:    router,
		Serv:      server,
	}
}

func (server *MicroHTTPServer) Initialize() http.Handler {
	return server.Router
}

func (server *MicroHTTPServer) AddEndpoint(endpoint string, f func(http.ResponseWriter, *http.Request)) {
	realEndpoint := server.APIPrefix + endpoint
	server.Router.HandleFunc(realEndpoint, f).Methods(http.MethodGet)
}
