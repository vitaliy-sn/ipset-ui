package main

import (
	"ipset-ui/internal/config"
	"ipset-ui/internal/ipset"
	"ipset-ui/internal/logger"
	"ipset-ui/internal/server"
)

func main() {
	config.Init()

	err := ipset.LoadAll()
	if err != nil {
		logger.Fatal(err.Error())
	}

	server.RunHTTPServer()
}
