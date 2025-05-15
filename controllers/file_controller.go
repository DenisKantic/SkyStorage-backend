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
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create directory"})
		return
	}

	var savedFiles []models.File

	for _, file := range files {
		dstPath := filepath.Join(uploadDir, fmt.Sprintf("%d_%s", time.Now().UnixNano(), file.Filename))
		if err := c.SaveUploadedFile(file, dstPath); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file"})
			return
		}

		fileModel := models.File{
			FileName: file.Filename,
			FileSize: int(file.Size),
			MimeType: file.Header.Get("Content-Type"),
			UploadAt: time.Now(),
			Path:     dstPath,
		}

		if err := database.DB.Create(&fileModel).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file"})
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
