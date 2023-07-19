package wsrpc

import (
	"net/http"
	"net/rpc"

	"golang.org/x/net/websocket"
)

// WSRPCServer is a websocket based rpc server.
type WSRPCServer struct {
	mux        *http.ServeMux
	httpServer *http.Server

	wsPath    string
	ws        *websocket.Conn
	wsHandler websocket.Handler

	RPCCodecFunc rpc.ServerCodec
	rpcServer    rpc.Server

	appName, token string
}

// NewServer returns a new WSRPCServer.
func NewServer(addr, wsPath string) *WSRPCServer {
	// create a http server for ws.
	mux := http.NewServeMux()
	httpServer := &http.Server{Addr: addr, Handler: mux}

	ss := &WSRPCServer{
		mux:        mux,
		httpServer: httpServer, // http server for ws
		wsPath:     wsPath,
		rpcServer:  *rpc.NewServer(), // rpc server
	}

	return ss
}

// Serve serves the server.
// It starts a http server and handle ws requests at path "/ws".
func (s *WSRPCServer) Serve() error {
	// handle ws request.
	s.mux.Handle(s.wsPath, s)

	return s.httpServer.ListenAndServe()
}

func (s *WSRPCServer) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	// upgrade to websockt
	wsHandler := websocket.Handler(s.serveConn)
	wsHandler.ServeHTTP(w, req)
}

// serve handles websocket requests from the peer.
// it starts a rpc server and handle this websocket connection as rpc connection.
func (s *WSRPCServer) serveConn(ws *websocket.Conn) {
	s.rpcServer.ServeConn(ws)
}

// Register register RPC services.
func (s *WSRPCServer) Register(service interface{}) {
	s.rpcServer.Register(service)
}

// Mux returns the http.ServeMux.
// You can use this mux to config more routers.
func (s *WSRPCServer) Mux() *http.ServeMux {
	return s.mux
}

// Close closes the server.
func (s *WSRPCServer) Close() error {
	return s.httpServer.Close()
}
