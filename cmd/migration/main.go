package main

import (
	"job_swipe/internal/database"
	"job_swipe/internal/models"
	"log"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, relying on environment variables")
	}

	database.Connect()

	if database.DB == nil {
		log.Fatal("Database not connected. Cannot run migrations.")
	}

	log.Println("Running database migrations...")
	err := database.DB.AutoMigrate(
		&models.User{},
		&models.Application{},
		&models.Company{},
		&models.Product{},
		&models.Message{},
		&models.JobSeekerProfile{},
		&models.Internship{},
		&models.Job{},
		&models.JobSwipe{},
		&models.JobProviderProfile{},
	)
	if err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}
	log.Println("Database migrations completed successfully.")
}
