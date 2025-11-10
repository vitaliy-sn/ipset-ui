package server

import (
	"fmt"
	"net/http/httputil"
	"net/url"
	"strings"

	"ipset-ui/internal/config"
	"ipset-ui/internal/controller"
	"ipset-ui/internal/logger"
	"ipset-ui/internal/static"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

var prefixesToCheck = []string{"/api/"}

// hasAnyPrefix проверяет, начинается ли путь с одного из указанных префиксов.
func hasAnyPrefix(path string, prefixes []string) bool {
	for _, prefix := range prefixes {
		if strings.HasPrefix(path, prefix) {
			return true
		}
	}
	return false
}

func NewRouter() *gin.Engine {
	r := gin.Default()

	// Middleware логгирования всех запросов
	r.Use(func(c *gin.Context) {
		logger.Info("HTTP Request",
			"method", c.Request.Method,
			"path", c.Request.URL.Path,
			"client_ip", c.ClientIP(),
		)
		c.Next()
	})

	// Setting up CORS
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	ipsetController := controller.NewIPSetController()

	api := r.Group("/api/v1")

	api.GET("/ipsets", ipsetController.GetSets)
	api.POST("/ipsets", ipsetController.CreateSet)
	api.DELETE("/ipsets/:setName", ipsetController.DestroySet)
	api.POST("/ipsets/:setName/entries", ipsetController.AddEntry)
	api.DELETE("/ipsets/:setName/entries", ipsetController.DeleteEntry)
	api.POST("/ipsets/:setName/entries/import", ipsetController.ImportEntries)
	api.POST("/ipsets/:setName/entries/search", ipsetController.GetEntries)
	api.GET("/ipsets/:setName/backups", ipsetController.ListBackupFiles)
	api.DELETE("/ipsets/:setName/backups", ipsetController.DeleteBackupFile)
	api.POST("/ipsets/:setName/save", ipsetController.CreateBackup)
	api.POST("/ipsets/:setName/restore", ipsetController.RestoreSet)
	api.POST("/whois", ipsetController.GetWhoisInfo)
	api.POST("/dns-lookup", ipsetController.DNSLookup)

	frontendURL := config.Config.FrontendURL

	// Use EmbedStaticHandler
	if frontendURL == "" {
		logger.Info("Static serving mode: embedded static handler")
		r.Use(func(c *gin.Context) {
			if hasAnyPrefix(c.Request.URL.Path, prefixesToCheck) {
				c.Next()
				return
			}
			static.EmbedStaticHandler().ServeHTTP(c.Writer, c.Request)
			c.Abort()
		})
	} else {
		// If FrontendURL is set, proxy all non-API requests to it
		remote, err := url.Parse(frontendURL)
		if err != nil {
			panic(fmt.Sprintf("Invalid FRONTEND_URL: %s", err))
		}

		proxy := httputil.NewSingleHostReverseProxy(remote)

		r.Use(func(c *gin.Context) {
			if !hasAnyPrefix(c.Request.URL.Path, prefixesToCheck) {
				proxy.ServeHTTP(c.Writer, c.Request)
				c.Abort()
				return
			}
			c.Next()
		})
	}

	return r
}
