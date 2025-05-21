package controllers

import "C"
import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
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
		originalFileName := file.Filename

		// Check if file with same name already exists in DB
		var existing models.File
		if err := database.DB.Where("file_name = ?", originalFileName).First(&existing).Error; err == nil {
			c.JSON(http.StatusConflict, gin.H{"error": fmt.Sprintf("File '%s' already exists", originalFileName)})
			return
		}

		dstPath := filepath.Join(uploadDir, originalFileName)

		if err := c.SaveUploadedFile(file, dstPath); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file"})
			return
		}

		fileModel := models.File{
			FileName: originalFileName,
			FileSize: int(file.Size),
			MimeType: file.Header.Get("Content-Type"),
			UploadAt: time.Now(),
			Path:     originalFileName, // only the filename
		}

		if err := database.DB.Create(&fileModel).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file info"})
			return
		}

		savedFiles = append(savedFiles, fileModel)
	}

	// invalidating cache key after saving in DB
	err = database.RedisClient.Del(database.Ctx, "all_files").Err()
	if err != nil {
		log.Println("Failed to invalidate cache key", err)
	}

	c.JSON(http.StatusOK, gin.H{
		"success": "Files uploaded successfully",
		"files":   savedFiles,
	})
}
func GetAllFiles(c *gin.Context) {
	cacheKey := "all_files"

	cached, err := database.RedisClient.Get(database.Ctx, cacheKey).Result()
	if err == nil {
		log.Println("Serving files from Redis cache")
		// cache hit, return cached JSON directly
		var files []models.File
		if err := json.Unmarshal([]byte(cached), &files); err == nil {
			// succesfully unmarshaled from cache, return cached response
			c.JSON(http.StatusOK, gin.H{"files": files})
			return
		}
		// if unmarshal fails, continue to fetch fresh data
	}

	// cache miss or unmarshal error, fetch from DB
	log.Println("Fetching files from database")

	var files []models.File
	if err := database.DB.Find(&files).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get the files"})
		return
	}

	// marshal files to JSON to cache
	jsonBytes, err := json.Marshal(files)
	if err == nil {
		// cache the results for 10 minutes
		database.RedisClient.Set(database.Ctx, cacheKey, jsonBytes, 10*time.Minute)
	}

	// return response
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
func DeleteFile(c *gin.Context) {
	var input struct {
		FileName string `json:"filename"`
	}

	if err := c.ShouldBindJSON(&input); err != nil || input.FileName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing or invalid filename"})
		return
	}

	// Check if the file exists in DB
	var file models.File
	if err := database.DB.Where("file_name = ?", input.FileName).First(&file).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "File not found in database"})
		return
	}

	// Build the full file path
	filePath := filepath.Join("./uploads", file.FileName)

	// Try to delete the physical file
	if err := os.Remove(filePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete file from disk"})
		return
	}

	// Delete DB record
	if err := database.DB.Delete(&file).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete file record from database"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": "File deleted successfully"})
}

// formatBytes converts bytes to a human-readable string (e.g., MB, GB)
func formatBytesInGB(bytes int64) string {
	gb := float64(bytes) / (1024 * 1024 * 1024)
	return fmt.Sprintf("%.2f GB", gb)
}
