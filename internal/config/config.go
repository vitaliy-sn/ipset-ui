package config

import (
	"fmt"
	"os"

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

	// Create AppDir and BackupDir if it doesn't exist
	if _, err := os.Stat(Config.BackupDir); os.IsNotExist(err) {
		err := os.MkdirAll(Config.BackupDir, 0755)
		if err != nil {
			fmt.Printf("Error creating AppDir: %v\n", err)
		} else {
			fmt.Printf("Created AppDir: %s\n", Config.BackupDir)
		}
	}

	fmt.Printf("Config initialized:\n")
	fmt.Printf("  ListenAddress: %s\n", Config.ListenAddress)
	fmt.Printf("  AppDir:        %s\n", Config.AppDir)
	fmt.Printf("  FrontendURL:   %s\n", Config.FrontendURL)
}
