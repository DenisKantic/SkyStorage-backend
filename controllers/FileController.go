package controllers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"path/filepath"
	"sky_storage_golang/database"
	"sky_storage_golang/models"
	"time"
)

func UploadFiles(c *gin.Context) {
	form, err := c.MultipartForm()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to parse the form"})
		return
	}

	files := form.File["files"]

	uploadDir := "./uploads"
	if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create upload directory"})
		return
	}

	var savedFiles []models.File

	for _, file := range files {
		// Generate a unique filename
		filename := fmt.Sprintf("%d_%s", time.Now().UnixNano(), file.Filename)

		// Relative path for saving to disk
		relativePath := filepath.Join("uploads", filename)
		// Absolute path for saving the file physically
		dstPath := filepath.Join(".", relativePath)

		if err := c.SaveUploadedFile(file, dstPath); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file"})
			return
		}

		// Store only the relative path in the database
		fileModel := models.File{
			FileName: file.Filename,
			FileSize: int(file.Size),
			MimeType: file.Header.Get("Content-Type"),
			UploadAt: time.Now(),
			Path:     relativePath, // store only "uploads/filename"
		}

		if err := database.DB.Create(&fileModel).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file info"})
			return
		}

		savedFiles = append(savedFiles, fileModel)
	}

	c.JSON(http.StatusOK, gin.H{"success": "Files uploaded successfully"})
}
func GetAllFiles(c *gin.Context) {
	var files []models.File

	if err := database.DB.Find(&files).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get the files"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"files": files})
}
func GetUploadsFolderSize(c *gin.Context) {
	folderPath := "./uploads"

	var totalSize int64

	err := filepath.Walk(folderPath, func(_ string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			totalSize += info.Size()
		}
		return nil
	})

	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to read uploads folder", "details": err.Error()})
		return
	}

	c.JSON(200, gin.H{
		"size_bytes": totalSize,
		"size_human": formatBytesInGB(totalSize),
	})
}

// formatBytes converts bytes to a human-readable string (e.g., MB, GB)
func formatBytesInGB(bytes int64) string {
	gb := float64(bytes) / (1024 * 1024 * 1024)
	return fmt.Sprintf("%.2f GB", gb)
}
