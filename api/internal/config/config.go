package config

import (
	"log"
	"os"

	"github.com/spf13/viper"
)

type Config struct {
	AppPort        string
	AppHost        string
	IpsetBackupDir string
	StaticDir      string
	FrontendURL    string
}

var AppConfig Config

func LoadConfig() {
	viper.AutomaticEnv()

	viper.SetDefault("APP_PORT", "8080")
	viper.SetDefault("APP_HOST", "0.0.0.0")

	AppConfig = Config{
		AppPort:        viper.GetString("APP_PORT"),
		AppHost:        viper.GetString("APP_HOST"),
		IpsetBackupDir: viper.GetString("IPSET_BACKUP_DIR"),
		StaticDir:      viper.GetString("STATIC_DIR"),
		FrontendURL:    viper.GetString("FRONTEND_URL"),
	}

	if AppConfig.IpsetBackupDir == "" {
		log.Fatalf("Environment variable IPSET_BACKUP_DIR must be set")
	}

	if _, err := os.Stat(AppConfig.IpsetBackupDir); os.IsNotExist(err) {
		log.Fatalf("Backup directory does not exist: %s", AppConfig.IpsetBackupDir)
	}
}
