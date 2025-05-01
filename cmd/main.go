package main

import (
	"github.com/gin-gonic/gin"
	"sky_storage_golang/routes"
)

func main() {

	r := gin.Default()

	routes.AuthRoutes(r)
	routes.LogoutRoute(r)
	routes.SendEmail(r)

	r.Run(":8080")
}
