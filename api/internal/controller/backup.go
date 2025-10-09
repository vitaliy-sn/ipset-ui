package controller

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// ListBackupFiles returns the list of backup files for the specified set.
func (c *IPSetController) ListBackupFiles(ctx *gin.Context) {
	setName := ctx.Param("setName")

	files, err := c.ipset.ListBackupFiles(setName)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, files)
}

// SaveSet saves the set to a backup file.
func (c *IPSetController) SaveSet(ctx *gin.Context) {
	var req struct {
		SetName      string `json:"setName"`
		FileNamePart string `json:"fileNamePart"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	backupFileName := fmt.Sprintf("%s-%s.save", req.SetName, req.FileNamePart)

	if strings.HasSuffix(req.FileNamePart, ".save") {
		backupFileName = req.FileNamePart
	}

	err := c.ipset.SaveSet(req.SetName, backupFileName)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Set saved"})
}

// RestoreSet restores the set from a backup file.
func (c *IPSetController) RestoreSet(ctx *gin.Context) {
	var req struct {
		SetName      string `json:"setName"`
		FileNamePart string `json:"fileNamePart"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := c.ipset.RestoreSet(req.SetName, req.FileNamePart)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Set restored"})
}

// DeleteBackupFile deletes the backup file for the specified set.
func (c *IPSetController) DeleteBackupFile(ctx *gin.Context) {
	setName := ctx.Param("setName")
	var req struct {
		FileNamePart string `json:"fileNamePart"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := c.ipset.DeleteBackupFile(setName, req.FileNamePart)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Backup file deleted"})
}
