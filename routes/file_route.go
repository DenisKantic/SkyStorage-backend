package routes

import (
	"github.com/gin-gonic/gin"
	"sky_storage_golang/controllers"
)

func UploadRoute(r *gin.Engine) {

	filesRoute := r.Group("/files")

	filesRoute.POST("/upload", controllers.UploadFiles)
	// get all files
	filesRoute.GET("/all-uploads", controllers.GetAllFiles)
	filesRoute.GET("/total-size", controllers.GetUploadsFolderSize)
}
