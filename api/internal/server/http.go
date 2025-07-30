package server

import (
	"fmt"
	"ipset-ui/internal/config"
	"ipset-ui/internal/controller"
	"net/http/httputil"
	"net/url"
	"os"
	"path/filepath"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// RunHTTPServer initializes and starts the gin HTTP server.
func RunHTTPServer() {
	config.LoadConfig()

	r := gin.Default()

	// Setting up CORS
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"}, // All domains are allowed. For production, specific domains should be specified
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	// If FrontendURL is set, proxy all non-API requests to it
	frontendURL := config.AppConfig.FrontendURL
	if frontendURL != "" {
		remote, err := url.Parse(frontendURL)
		if err != nil {
			panic(fmt.Sprintf("Invalid FRONTEND_URL: %s", err))
		}

		proxy := httputil.NewSingleHostReverseProxy(remote)

		r.Use(func(c *gin.Context) {
			if !filepath.HasPrefix(c.Request.URL.Path, "/api/") {
				proxy.ServeHTTP(c.Writer, c.Request)
				c.Abort()
				return
			}
			c.Next()
		})
	} else {
		// If FRONTEND_URL is not set, use static files from STATIC_DIR if it exists.
		// Only non-API requests (not starting with /api/) will be handled by this middleware.
		staticDirExists := true
		if _, err := os.Stat(config.AppConfig.StaticDir); os.IsNotExist(err) {
			staticDirExists = false
		}

		if staticDirExists {
			// Middleware for serving static files
			r.Use(func(c *gin.Context) {
				if !filepath.HasPrefix(c.Request.URL.Path, "/api/") {
					filePath := filepath.Join(config.AppConfig.StaticDir, c.Request.URL.Path)
					if _, err := os.Stat(filePath); os.IsNotExist(err) {
						// If the file does not exist, return index.html
						c.File(filepath.Join(config.AppConfig.StaticDir, "index.html"))
					} else {
						// If the file exists, return it
						c.File(filePath)
					}
					c.Abort()
					return
				}
				c.Next()
			})
		}
	}

	ipsetController := controller.NewIPSetController()

	api := r.Group("/api/v1")
	{
		api.GET("/ipsets", ipsetController.GetSets)
		api.POST("/ipsets", ipsetController.CreateSet)
		api.DELETE("/ipsets/:setName", ipsetController.DeleteSet)
		api.POST("/ipsets/:setName/entries", ipsetController.AddEntry)
		api.DELETE("/ipsets/:setName/entries", ipsetController.DeleteEntry)
		api.POST("/ipsets/:setName/entries/import", ipsetController.ImportEntries)
		api.POST("/ipsets/:setName/entries/search", ipsetController.GetEntries)
		api.GET("/ipsets/:setName/backups", ipsetController.ListBackupFiles)
		api.DELETE("/ipsets/:setName/backups", ipsetController.DeleteBackupFile)
		api.POST("/ipsets/:setName/save", ipsetController.SaveSet)
		api.POST("/ipsets/:setName/restore", ipsetController.RestoreSet)
		api.POST("/whois", ipsetController.GetWhoisInfo)
		api.POST("/dns-lookup", ipsetController.DNSLookup)
	}

	address := fmt.Sprintf("%s:%s", config.AppConfig.AppHost, config.AppConfig.AppPort)
	r.Run(address)
}
