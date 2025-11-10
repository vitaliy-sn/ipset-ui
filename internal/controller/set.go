package controller

import (
	"ipset-ui/internal/ipset"
	"net/http"

	"github.com/gin-gonic/gin"
)

// CreateSet handles the creation of a new set.
func (c *IPSetController) CreateSet(ctx *gin.Context) {
	var req struct {
		SetName string `json:"setName"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := ipset.Create(req.SetName); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Set created"})
}

// DestroySet handles the deletion of a set.
func (c *IPSetController) DestroySet(ctx *gin.Context) {
	setName := ctx.Param("setName")

	if err := ipset.Destroy(setName); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Set deleted"})
}

// GetSets handles listing all sets.
func (c *IPSetController) GetSets(ctx *gin.Context) {
	sets, err := ipset.List()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, sets)
}
