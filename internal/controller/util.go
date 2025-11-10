package controller

import (
	"ipset-ui/internal/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetWhoisInfo performs a whois lookup for the specified object.
func (c *IPSetController) GetWhoisInfo(ctx *gin.Context) {
	var req struct {
		Object string `json:"object"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	info, err := utils.Whois(req.Object)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.String(http.StatusOK, info)
}

// DNSLookup performs a DNS lookup for the specified domain.
func (c *IPSetController) DNSLookup(ctx *gin.Context) {
	var req struct {
		Domain string `json:"domain"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ips, err := utils.LookupIPv4(req.Domain)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, ips)
}
