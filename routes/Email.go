package routes

import (
	"github.com/gin-gonic/gin"
	"sky_storage_golang/controllers"
)

func SendEmail(r *gin.Engine) {
	emailGroup := r.Group("/email")

	emailGroup.POST("/send-email", controllers.SendEmail)
}
