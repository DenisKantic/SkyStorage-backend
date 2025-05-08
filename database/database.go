package database

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"os"
	"sky_storage_golang/models"
)

var DB *gorm.DB

func ConnectDB() {
	// Use environment variables for security
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
	)

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	err = DB.AutoMigrate(&models.User{}, &models.Email{})
	if err != nil {
		log.Fatal("Failed to migrate table:", err)
	}
	//Create a default dev user if not exists
	//var devUser models.User
	//result := DB.First(&devUser, "username = ?", "devadmin")
	//if result.Error != nil {
	//	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
	//		hashedPassword, hashErr := bcrypt.GenerateFromPassword([]byte("devpass123"), bcrypt.DefaultCost)
	//		if hashErr != nil {
	//			log.Fatalf("Failed to hash password: %v", hashErr)
	//		}
	//
	//		// Hash password before saving in real-world apps!
	//		defaultUser := models.User{
	//			Username: "denis",
	//			Password: string(hashedPassword), // ⚠️ In production, hash this!
	//			Email:    "denis.kantic18@gmail.com",
	//		}
	//		if err := DB.Create(&defaultUser).Error; err != nil {
	//			log.Fatalf("Failed to create default user: %v", err)
	//		} else {
	//			log.Println("Default dev user created")
	//		}
	//	} else {
	//		log.Fatalf("Error checking for default user: %v", result.Error)
	//	}
	//}

	log.Println("Migration completed")
}
