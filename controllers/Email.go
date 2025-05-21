package controllers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"html"
	"net/http"
	"net/smtp"
	"os"
	"sky_storage_golang/database"
	"sky_storage_golang/models"
	"time"
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

	date := time.Now().Format(time.RFC1123Z)
	messageID := fmt.Sprintf("%d@vortexdigitalsystems.com", time.Now().UnixNano())

	// HTML content with footer
	htmlBody := fmt.Sprintf(`
		<!DOCTYPE html>
		<html>
		<body style="font-family: sans-serif; color: #333;">
			<div>
				<pre style="white-space: pre-wrap; font-family: inherit; font-size: 14px;">%s</pre>
				<hr style="border: none; border-top: 1px solid #ccc; margin-top: 30px;">
				<p style="font-size: 12px; color: #555;">
					<strong>Denis Kantic</strong><br>
					Vortex Digital Systems<br>
					Full Stack Web Developer &amp; SysAdmin <br><br>
					Email sent: %s
				</p>
			</div>
		</body>
		</html>
	`, html.EscapeString(email.Body), date)

	msg := []byte("From: \"Denis Kantic\" <" + smtpUser + ">\r\n" +
		"To: <" + email.To + ">\r\n" +
		"Subject: " + email.Subject + "\r\n" +
		"Date: " + date + "\r\n" +
		"Message-ID: <" + messageID + ">\r\n" +
		"MIME-Version: 1.0\r\n" +
		"Content-Type: text/html; charset=\"UTF-8\"\r\n" +
		"\r\n" +
		htmlBody + "\r\n")

	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, smtpUser, []string{email.To}, msg)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to send email"})
		return
	}

	email.SentAt = time.Now() // getting the time for this email when sent

	// save sent email to DB
	if err := database.DB.Create(&email).Error; err != nil {
		c.JSON(500, gin.H{"error": "Failed to save email"})
		return
	}
	c.JSON(200, gin.H{"success": "Email sent successfully"})
}

func ServeSentEmails(c *gin.Context) {
	var emails []models.Email

	fmt.Println("ServeSentEmails triggered") // Debug line

	if err := database.DB.Order("sent_at DESC").Find(&emails).Error; err != nil {
		c.JSON(500, gin.H{"error": "Failed to read sent emails"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"emails": emails})
}

func SendLoginCode(to, subject, body string) error {
	smtpHost := os.Getenv("SMTP_HOST")
	smtpPort := os.Getenv("SMTP_PORT")
	smtpUser := os.Getenv("SMTP_USER")
	smtpPassword := os.Getenv("SMTP_PASSWORD")

	auth := smtp.PlainAuth("", smtpUser, smtpPassword, smtpHost)

	date := time.Now().Format(time.RFC1123Z)
	messageID := fmt.Sprintf("%d@vortexdigitalsystems.com", time.Now().UnixNano())

	msg := []byte("From: " + smtpUser + "\r\n" +
		"To: " + to + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"Date: " + date + "\r\n" +
		"Message-ID: <" + messageID + ">\r\n" +
		"MIME-Version: 1.0\r\n" +
		"Content-Type: text/plain; charset=\"UTF-8\"\r\n" +
		"\r\n" +
		body + "\r\n")

	return smtp.SendMail(smtpHost+":"+smtpPort, auth, smtpUser, []string{to}, msg)
}
