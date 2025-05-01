package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"sky_storage_golang/config"
	"sky_storage_golang/routes"
)

func main() {

	gin.SetMode(gin.ReleaseMode)

	cors_config := cors.Config{
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization", "X-CSRF-Token"},
		AllowCredentials: true,
		AllowOrigins:     []string{"http://localhost:5173"},
	}

	r := gin.Default()
	r.Use(cors.New(cors_config))
	config.LoadEnv()

	routes.AuthRoutes(r)
	routes.LogoutRoute(r)
	routes.SendEmail(r)

	r.Run(":8080")
}
