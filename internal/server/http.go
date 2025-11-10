package server

import (
	"ipset-ui/internal/config"
)

// RunHTTPServer initializes and starts the gin HTTP server.
func RunHTTPServer() {
	r := NewRouter()
	r.Run(config.Config.ListenAddress)
}
