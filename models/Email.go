package models

import "time"

type Email struct {
	ID      uint      `gorm:"primary_key" json:"id"`
	To      string    `json:"to" binding:"required"`
	Subject string    `json:"subject" binding:"required"`
	Body    string    `json:"body" binding:"required"`
	SentAt  time.Time `json:"sent_at"`
}
