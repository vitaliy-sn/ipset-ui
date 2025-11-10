package controller

import (
	"fmt"
	"ipset-ui/internal/config"
	"ipset-ui/internal/ipset"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
)

// ListBackupFiles returns the list of backup files for the specified set.
func (c *IPSetController) ListBackupFiles(ctx *gin.Context) {
	setName := ctx.Param("setName")

	files, err := ipset.ListBackups(setName)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("failed to list backup files for set '%s': %v", setName, err),
		})
		return
	}

	ctx.JSON(http.StatusOK, files)
}

// SaveSet saves the set to a backup file.
func (c *IPSetController) CreateBackup(ctx *gin.Context) {
	var req struct {
		SetName      string `json:"setName"`
		FileNamePart string `json:"fileNamePart"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("invalid request format: %v", err),
		})
		return
	}

	backupFileName := fmt.Sprintf("%s-%s.save", req.SetName, req.FileNamePart)

	if strings.HasSuffix(req.FileNamePart, ".save") {
		backupFileName = req.FileNamePart
	}

	err := ipset.CreateBackup(req.SetName, backupFileName)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("failed to create backup for set '%s' as file '%s': %v", req.SetName, backupFileName, err),
		})
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
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("invalid request format: %v", err),
		})
		return
	}

	err := ipset.Flush(req.SetName)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("failed to flush set '%s': %v", req.SetName, err),
		})
		return
	}

	fileName := fmt.Sprintf("%s-%s.save", req.SetName, req.FileNamePart)
	path := filepath.Join(config.Config.BackupDir, fileName)

	err = ipset.Restore(path)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("failed to restore set '%s' from file '%s': %v", req.SetName, fileName, err),
		})
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
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("invalid request format: %v", err),
		})
		return
	}

	err := ipset.DeleteBackup(setName, req.FileNamePart)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("failed to delete backup file '%s' for set '%s': %v", req.FileNamePart, setName, err),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Backup file deleted"})
}
