package models

type User struct {
	ID       uint   `gorm:"primary_key"`
	Username string `gorm:"not null;unique" json:"username" binding:"required"`
	Password string `gorm:"not null" json:"password" binding:"required"`
	Email    string `gorm:"not null;unique" json:"email" binding:"required"`
}
