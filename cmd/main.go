package main

import (
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"sky_storage_golang/config"
	"sky_storage_golang/database"
	"sky_storage_golang/routes"
)

func main() {

	gin.SetMode(gin.ReleaseMode)

	cors_config := cors.Config{
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization", "X-CSRF-Token"},
		AllowCredentials: true,
		//AllowAllOrigins:  true,
		AllowOrigins: []string{"http://localhost:5173"},
	}

	r := gin.Default()
	r.Use(cors.New(cors_config))
	config.LoadEnv()
	database.ConnectDB()

	routes.AuthRoutes(r)
	routes.LogoutRoute(r)
	routes.EmailRoutes(r)
	routes.UploadRoute(r)

	// SERVING FOLDER
	r.Static("/uploads", "./uploads")

	err := r.Run(":8080")
	if err != nil {
		fmt.Println("Error starting at 8080 port")
		return
	}
}
