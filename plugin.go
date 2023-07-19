package wsrpc

import (
	"net/http"

	"golang.org/x/net/websocket"
)

// server plugins

// ServeHTTPFunc is invoked before the http connection is made up.
type ServeHTTPFunc func(w http.ResponseWriter, req *http.Request) error

// ServeWSConnFunc is invoked before the websocket connection is made up.
type ServeWSConnFunc func(ws *websocket.Conn) error

// RegisterFunc is invoked when a service is registered.
type RegisterFunc func(service interface{}) error

// client plugins

// DialFunc is invoked before the websocket connection is made up.
type DialFunc func(config *websocket.Config) error

// BeforeCallFunc is invoked before the rpc call.
type BeforeCallFunc func(serviceMethod string, args interface{}) error

// AfterCallFunc is invoked after the rpc call.
type AfterCallFunc func(serviceMethod string, args interface{}, reply interface{}, err error) error

// BeforeGoFunc is invoked before the rpc call.
type BeforeGoFunc func(serviceMethod string, args interface{}) error
