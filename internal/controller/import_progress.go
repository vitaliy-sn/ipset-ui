package controller

import (
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
)

// ImportProgress tracks the state of the currently running file import.
// Only one import can run at a time: the UI disables the import button
// while it is in progress.
type ImportProgress struct {
	SetName   string `json:"setName"`
	Total     int    `json:"total"`
	Processed int    `json:"processed"`
	Added     int    `json:"added"`
	Done      bool   `json:"done"`
	Error     string `json:"error,omitempty"`
}

var (
	importProgressMu sync.Mutex
	currentImport    *ImportProgress
)

func setImportProgress(p *ImportProgress) {
	importProgressMu.Lock()
	defer importProgressMu.Unlock()
	currentImport = p
}

func updateImportProgress(added, processed int) {
	importProgressMu.Lock()
	defer importProgressMu.Unlock()
	if currentImport == nil {
		return
	}
	currentImport.Processed = processed
	currentImport.Added = added
}

func finishImportProgress(added, total int, errorMessage string) {
	importProgressMu.Lock()
	defer importProgressMu.Unlock()
	if currentImport == nil {
		return
	}
	currentImport.Done = true
	currentImport.Processed = total
	currentImport.Added = added
	currentImport.Error = errorMessage
}

func getImportProgress() *ImportProgress {
	importProgressMu.Lock()
	defer importProgressMu.Unlock()
	if currentImport == nil {
		return nil
	}
	progress := *currentImport
	return &progress
}

// GetImportProgress returns the progress of the currently running file import.
func (c *IPSetController) GetImportProgress(ctx *gin.Context) {
	progress := getImportProgress()
	if progress == nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "no import in progress"})
		return
	}
	ctx.JSON(http.StatusOK, progress)
}
