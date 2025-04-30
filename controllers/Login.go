package controllers

import (
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"sky_storage_golang/database"
	"sky_storage_golang/models"
	"sky_storage_golang/utils"
)

func Login(c *gin.Context) {
	var input struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		fmt.Printf("ERROR:", err.Error())
		return
	}

	var user models.UserEmployee
	if err := database.DB.Where("username = ?", input.Username).First(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User not found"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid password"})
		return
	}

	token, err := utils.GenerateJWT(user.ID, user.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to login"})
		return
	}

	c.SetCookie("auth_token", token, 60*60*24, "/", "localhost", false, false)
	c.Header("Set-Cookie", "auth_token="+token+"; Max-Age=86400; Path=/; Domain=localhost; SameSite=Strict; HttpOnly")
	c.JSON(http.StatusOK, gin.H{"success": "Logged in"})
}

func Logout(c *gin.Context) {
	c.SetCookie("auth_token", "", -1, "/", "localhost", false, false)

	session := sessions.Default(c)
	session.Clear()
	session.Save()

	c.SetCookie("my-session", "", -1, "/", "localhost", false, false)

	c.JSON(200, gin.H{"success": "Logout successfully"})
}
