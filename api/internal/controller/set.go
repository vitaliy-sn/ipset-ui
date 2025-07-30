package controller

import (
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

	exists, err := c.ipset.SetExists(req.SetName)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if exists {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "set already exists"})
		return
	}

	if err := c.ipset.CreateSet(req.SetName); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Set created"})
}

// DeleteSet handles the deletion of a set.
func (c *IPSetController) DeleteSet(ctx *gin.Context) {
	setName := ctx.Param("setName")

	exists, err := c.ipset.SetExists(setName)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if !exists {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "set does not exist"})
		return
	}

	if err := c.ipset.DeleteSet(setName); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Set deleted"})
}

// GetSets handles listing all sets.
func (c *IPSetController) GetSets(ctx *gin.Context) {
	sets, err := c.ipset.ListSets()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, sets)
}
