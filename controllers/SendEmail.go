package controllers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/smtp"
	"os"
	"sky_storage_golang/models"
)

func SendEmail(c *gin.Context) {
	var email models.Email

	if err := c.ShouldBindJSON(&email); err != nil {
		c.JSON(400, gin.H{"error": "Some fields are missing"})
		fmt.Println(err.Error())
		return
	}

	smtpHost := os.Getenv("SMTP_HOST")
	smtpPort := os.Getenv("SMTP_PORT")
	smtpUser := os.Getenv("SMTP_USER")
	smtpPassword := os.Getenv("SMTP_PASSWORD")

	auth := smtp.PlainAuth("", smtpUser, smtpPassword, smtpHost)

	msg := []byte("From: \"" + "Denis Kantic" + "\" <" + smtpUser + ">\r\n" +
		"To: \"" + email.To + "\"\r\n" +
		"Subject: " + email.Subject + "\r\n" +
		"Content-Type: text/plain; charset=\"UTF-8\"\r\n" +
		"\r\n" +
		email.Body + "\r\n")

	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, smtpUser, []string{email.To}, msg)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to send email"})
		return
	}

	c.JSON(200, gin.H{"success": "Email sent successfully"})
}
