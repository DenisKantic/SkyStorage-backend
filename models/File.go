package models

import "time"

type File struct {
	ID       uint      `gorm:"primary_key;auto_increment" json:"id"`
	FileName string    `gorm:"not null" json:"file_name"`
	FileSize int       `gorm:"not null" json:"file_size"`
	MimeType string    `gorm:"not null" json:"mime_type"`
	UploadAt time.Time `gorm:"autoCreateTime" json:"upload_at"`
	Path     string    `gorm:"not null" json:"path"`
}
