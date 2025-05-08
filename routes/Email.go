package routes

import (
	"github.com/gin-gonic/gin"
	"sky_storage_golang/controllers"
)

func EmailRoutes(r *gin.Engine) {
	emailGroup := r.Group("/email")

	emailGroup.POST("/send-email", controllers.SendEmail)
	emailGroup.GET("/sent-emails", controllers.ServeSentEmails)
}
