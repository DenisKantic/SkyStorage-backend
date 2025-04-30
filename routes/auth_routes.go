package routes

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	csrf "github.com/utrack/gin-csrf"
	"sky_storage_golang/controllers"
	"sky_storage_golang/middleware"
)

func AuthRoutes(r *gin.Engine) {

	protectedGroup := r.Group("/protected")
	protectedGroup.Use(middleware.AuthMiddleware())

	// login user group (login, register)
	userGroup := r.Group("/auth")
	{
		userGroup.POST("/login", controllers.Login)
	}

	// protected routes (need to be logged in with valid JWT token)
	protectedGroup.GET("/profile", func(c *gin.Context) {
		c.JSON(200, gin.H{"success": "working"})
	})
}

func LogoutRoute(r *gin.Engine) {
	store := cookie.NewStore([]byte("super-secret-key"))
	r.Use(sessions.Sessions("my-session", store)) // <- Important

	csrfGroup := r.Group("/")
	csrfGroup.Use(csrf.Middleware(csrf.Options{
		Secret: "very-secret-csrf-key",
		ErrorFunc: func(c *gin.Context) {
			c.JSON(403, gin.H{"error": "CSRF token invalid or missing"})
			c.Abort()
		},
	}))

	// Provide CSRF token to frontend (GET request)
	csrfGroup.GET("/csrf-token", func(c *gin.Context) {
		c.JSON(200, gin.H{"csrf_token": csrf.GetToken(c)})
	})

	// Protected logout
	csrfGroup.POST("/auth/logout", controllers.Logout)

}
