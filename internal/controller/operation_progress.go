package controller

import (
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
)

// OperationProgress tracks the state of the currently running entries
// operation (file import or bulk delete).
// Only one operation can run at a time: the UI disables the corresponding
// buttons while it is in progress.
type OperationProgress struct {
	Operation string `json:"operation"` // "import" or "delete"
	SetName   string `json:"setName"`
	Total     int    `json:"total"`
	Processed int    `json:"processed"`
	Changed   int    `json:"changed"` // entries added (import) or deleted
	Done      bool   `json:"done"`
	Error     string `json:"error,omitempty"`
}

var (
	operationProgressMu sync.Mutex
	currentOperation    *OperationProgress
)

func setOperationProgress(p *OperationProgress) {
	operationProgressMu.Lock()
	defer operationProgressMu.Unlock()
	currentOperation = p
}

func updateOperationProgress(changed, processed int) {
	operationProgressMu.Lock()
	defer operationProgressMu.Unlock()
	if currentOperation == nil {
		return
	}
	currentOperation.Processed = processed
	currentOperation.Changed = changed
}

func finishOperationProgress(changed, total int, errorMessage string) {
	operationProgressMu.Lock()
	defer operationProgressMu.Unlock()
	if currentOperation == nil {
		return
	}
	currentOperation.Done = true
	currentOperation.Processed = total
	currentOperation.Changed = changed
	currentOperation.Error = errorMessage
}

func getOperationProgress() *OperationProgress {
	operationProgressMu.Lock()
	defer operationProgressMu.Unlock()
	if currentOperation == nil {
		return nil
	}
	progress := *currentOperation
	return &progress
}

// GetOperationProgress returns the progress of the currently running entries
// operation for the given set.
func (c *IPSetController) GetOperationProgress(ctx *gin.Context) {
	progress := getOperationProgress()
	if progress == nil || progress.SetName != ctx.Param("setName") {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "no operation in progress"})
		return
	}
	ctx.JSON(http.StatusOK, progress)
}
