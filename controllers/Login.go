package controllers

import (
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"math/rand"
	"net/http"
	"sky_storage_golang/database"
	"sky_storage_golang/models"
	"sky_storage_golang/utils"
	"time"
)

func Login(c *gin.Context) {

	userIP := c.ClientIP()

	var input struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		fmt.Printf("ERROR:", err.Error())
		return
	}

	var user models.User
	if err := database.DB.Where("username = ?", input.Username).First(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User not found"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid password"})
		fmt.Println(err.Error())
		return
	}

	// check if IP is already verified
	var verificaton models.LoginVerification
	err := database.DB.Where("user_id = ? AND ip_address = ? AND verified = true", user.ID, userIP).First(&verificaton).Error

	if err == nil {
		// IP is verified - generate token and login
		token, err := utils.GenerateJWT(user.ID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to login"})
			return
		}

		c.SetCookie("auth_token", token, 60*60*24, "/", "localhost", false, false)
		c.JSON(http.StatusOK, gin.H{"success": "Logged in"})
		return
	}

	// IP is not verified - generate login code
	code := GenerateSixDigitCode()
	expiresAt := time.Now().Add(time.Minute * 5)

	database.DB.Where("user_id = ? AND ip_address = ?", user.ID, userIP).Delete(&models.LoginVerification{})
	verificaton = models.LoginVerification{
		UserID:    user.ID,
		IPAddress: userIP,
		Code:      code,
		ExpiresAt: expiresAt,
		Verified:  false,
	}

	database.DB.Create(&verificaton)

	body := fmt.Sprintf(`We noticed a suspicious login attempt to your web app from IP address: %s

If this was you, please enter the following code to continue: %s

If you did not attempt to log in, we recommend resetting your password immediately.`, userIP, code)

	// send code via email
	err = SendLoginCode(user.Email, "Suspicious Login - Verification Code", body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send verification email"})
		return
	}

	// respond that verification code is sent and login is pending
	c.JSON(http.StatusAccepted, gin.H{
		"error":         "Verification code sent. Please check your email.",
		"code_required": true,
	})
}

func GenerateSixDigitCode() string {
	// Create a new source and random generator
	source := rand.NewSource(time.Now().UnixNano())
	r := rand.New(source)

	// Generate zero-padded 6-digit number
	code := fmt.Sprintf("%06d", r.Intn(1000000))
	return code
}

func Logout(c *gin.Context) {
	c.SetCookie("auth_token", "", -1, "/", "localhost", false, false)

	session := sessions.Default(c)
	session.Clear()
	session.Save()

	c.SetCookie("my-session", "", -1, "/", "localhost", false, false)

	c.JSON(200, gin.H{"success": "Logout successfully"})
}
