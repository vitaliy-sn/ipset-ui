package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type config struct {
	ListenAddress string
	AppDir        string
	FrontendURL   string
	BackupDir     string
}

var Config config

func Init() {
	viper.AutomaticEnv()

	viper.SetDefault("LISTEN_ADDRESS", "0.0.0.0:8080")
	viper.SetDefault("APP_DIR", "/opt/ipset-ui")

	Config = config{
		ListenAddress: viper.GetString("LISTEN_ADDRESS"),
		FrontendURL:   viper.GetString("FRONTEND_URL"),
		AppDir:        viper.GetString("APP_DIR"),
	}

	Config.BackupDir = Config.AppDir + "/backups"

	fmt.Printf("Config initialized:\n")
	fmt.Printf("  ListenAddress: %s\n", Config.ListenAddress)
	fmt.Printf("  AppDir:        %s\n", Config.AppDir)
	fmt.Printf("  FrontendURL:   %s\n", Config.FrontendURL)
}
