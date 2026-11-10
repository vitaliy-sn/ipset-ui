package controller

import (
	"fmt"
	"ipset-ui/internal/ipset"
	"ipset-ui/internal/logger"
	"ipset-ui/internal/utils"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// AddEntry adds an entry to the set.
func (c *IPSetController) AddEntry(ctx *gin.Context) {
	setName := ctx.Param("setName")
	var req struct {
		Entry   string `json:"entry"`
		Comment string `json:"comment"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	entries := utils.SplitEntries(req.Entry)
	if len(entries) == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "no entries provided"})
		return
	}

	added, failed, err := ipset.AddEntries(setName, entries, req.Comment, nil)
	if err != nil {
		if added == 0 {
			if strings.Contains(err.Error(), "entry not added") {
				ctx.JSON(http.StatusConflict, gin.H{"error": "Entry already exists or not added. " + err.Error()})
				return
			}
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{
			"message": fmt.Sprintf("Added %d of %d entries", added, len(entries)),
			"added":   added,
			"total":   len(entries),
			"failed":  failed,
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Entry added",
		"added":   added,
		"total":   len(entries),
	})
}

// DeleteEntry removes an entry from the set.
func (c *IPSetController) DeleteEntry(ctx *gin.Context) {
	setName := ctx.Param("setName")
	var req struct {
		Entry string `json:"entry"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := ipset.DeleteEntry(setName, req.Entry); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Entry deleted"})
}

// DeleteEntries removes all entries matching the given filter.
func (c *IPSetController) DeleteEntries(ctx *gin.Context) {
	setName := ctx.Param("setName")
	var req struct {
		Filter string `json:"filter"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if strings.TrimSpace(req.Filter) == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "filter is required"})
		return
	}

	entries, err := ipset.ListEntries(setName)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	matched := ipset.FilterEntries(entries, req.Filter)
	if len(matched) == 0 {
		ctx.JSON(http.StatusOK, gin.H{"message": "No entries to delete"})
		return
	}

	names := make([]string, 0, len(matched))
	for _, entry := range matched {
		names = append(names, entry.Entry)
	}

	deleted, failed, err := ipset.DeleteEntries(setName, names)
	if err != nil && deleted == 0 {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if len(failed) > 0 {
		ctx.JSON(http.StatusOK, gin.H{
			"message": fmt.Sprintf("Deleted %d of %d entries", deleted, len(matched)),
			"deleted": deleted,
			"total":   len(matched),
			"failed":  failed,
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf("Deleted %d entries", deleted),
		"deleted": deleted,
		"total":   len(matched),
	})
}

// ImportEntries imports entries from an uploaded file (multipart/form-data).
// The import runs synchronously in this request; its progress is published
// in the shared import state and can be polled via GetImportProgress.
func (c *IPSetController) ImportEntries(ctx *gin.Context) {
	setName := ctx.Param("setName")

	exists, err := ipset.SetExists(setName)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list sets: " + err.Error()})
		return
	}
	if !exists {
		ctx.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprintf("set '%s' not found", setName)})
		return
	}

	if err := ctx.Request.ParseMultipartForm(10 << 20); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "failed to parse multipart form"})
		return
	}

	file, header, err := ctx.Request.FormFile("file")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "file is required"})
		return
	}
	defer file.Close()

	comment := ctx.Request.FormValue("comment")

	entries, err := utils.ReadEntriesFromReader(file)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "failed to read entries from file"})
		return
	}

	if len(entries) == 0 {
		ctx.JSON(http.StatusOK, gin.H{"message": "No entries to import", "added": 0})
		return
	}

	total := len(entries)
	setImportProgress(&ImportProgress{SetName: setName, Total: total})

	added, failed, err := ipset.AddEntries(setName, entries, comment, func(added, processed, total int) {
		updateImportProgress(added, processed)
	})

	errMessage := ""
	if err != nil {
		errMessage = fmt.Sprintf("Import from '%s' failed: %s", header.Filename, truncateMessage(err.Error()))
		logger.Warn("Import from file failed",
			"set", setName,
			"file", header.Filename,
			"added", added,
			"failed", len(failed),
			"error", errMessage,
		)
	}

	finishImportProgress(added, total, errMessage)

	if err != nil {
		// Partial success is reported as a completed import with failures;
		// otherwise (nothing added, or saving the set failed) it is an error.
		if added > 0 && len(failed) > 0 {
			ctx.JSON(http.StatusOK, gin.H{
				"message": fmt.Sprintf("Imported %d of %d entries", added, total),
				"added":   added,
				"failed":  len(failed),
				"total":   total,
			})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": errMessage})
		return
	}

	logger.Info("Import from file completed",
		"set", setName,
		"file", header.Filename,
		"added", added,
	)
	ctx.JSON(http.StatusOK, gin.H{"message": "Entries imported", "added": added})
}

// truncateMessage trims a long error message to keep the response compact.
func truncateMessage(s string) string {
	const limit = 200
	if len(s) <= limit {
		return s
	}
	return s[:limit] + "..."
}

// GetEntries returns the entries of the set, optionally with filtering.
func (c *IPSetController) GetEntries(ctx *gin.Context) {
	setName := ctx.Param("setName")

	var req struct {
		Filter string `json:"filter"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil && err.Error() != "EOF" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	entries, err := ipset.ListEntries(setName)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if req.Filter != "" {
		entries = ipset.FilterEntries(entries, req.Filter)
	}

	if entries == nil {
		entries = []ipset.EntryWithComment{}
	}

	ctx.JSON(http.StatusOK, entries)
}
