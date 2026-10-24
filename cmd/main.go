package main

import (
	"flag"
	"fmt"

	"ipset-ui/internal/config"
	"ipset-ui/internal/ipset"
	"ipset-ui/internal/logger"
	"ipset-ui/internal/server"
	"ipset-ui/internal/version"
)

func main() {
	showVersion := flag.Bool("version", false, "show build information and exit")
	flag.Parse()

	if *showVersion {
		fmt.Printf("ipset-ui (build: %s, commit: %s)\n", version.BuildTime, version.Commit)
		return
	}

	logger.Info("ipset-ui started", "build_time", version.BuildTime, "commit", version.Commit)

	config.Init()

	err := ipset.LoadAll()
	if err != nil {
		logger.Fatal(err.Error())
	}

	server.RunHTTPServer()
}
