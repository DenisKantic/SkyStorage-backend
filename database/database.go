package database

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"os"
	"sky_storage_golang/models"
)

// redis and postgres
var (
	DB          *gorm.DB
	RedisClient *redis.Client
	Ctx         = context.Background()
)

func ConnectRedis() {
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_ADDR"),     // e.g. "localhost:6379"
		Password: os.Getenv("REDIS_PASSWORD"), // empty if no password
		DB:       0,                           // default DB
	})

	// Test connection
	_, err := RedisClient.Ping(Ctx).Result()
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}

	log.Println("Connected to Redis")
}

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

	err = DB.AutoMigrate(&models.User{}, &models.Email{}, &models.File{}, &models.LoginVerification{})
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

	ConnectRedis()
}
