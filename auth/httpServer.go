package auth

import (
	"context"
	"net/http"
)

// HttpServer wraps the http.Server type for enhanced testability.
type HttpServer interface {
	// ListenAndServe listens on the TCP network address and then calls Serve to handle requests on
	// incoming connections.
	ListenAndServe() error
	// Shutdown shuts down the server without interrupting any active connections.
	Shutdown(ctx context.Context) error
}

// NewHttpServer creates a new instance of HttpServer.
func NewHttpServer(addr string) HttpServer {
	return &http.Server{Addr: addr}
}
