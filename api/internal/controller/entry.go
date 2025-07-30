package controller

import (
	"fmt"
	"ipset-ui/internal/ipset"
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

	if err := c.ipset.AddEntry(setName, req.Entry, req.Comment); err != nil {
		if strings.Contains(err.Error(), "entry not added") {
			ctx.JSON(http.StatusConflict, gin.H{"error": "Entry already exists or not added"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Entry added"})
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

	if err := c.ipset.DeleteEntry(setName, req.Entry); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Entry deleted"})
}

// ImportEntries imports entries from an uploaded file (multipart/form-data).
func (c *IPSetController) ImportEntries(ctx *gin.Context) {
	setName := ctx.Param("setName")

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

	added, err := c.ipset.AddEntries(setName, entries, comment)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("Import from '%s' failed: %s", header.Filename, err.Error()),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Entries imported", "added": added})
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

	entries, err := c.ipset.ListEntries(setName)
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
