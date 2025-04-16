package initializers

import (
	"fmt"
	"go-stripe-chi-example/models"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Global DB variable
var DB *gorm.DB

// Initialize Database
func InitDatabase() {
	dsn := os.Getenv("DSN")
	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{}) // Assign to global DB
	if err != nil {
		log.Fatalf("Cannot connect to database: %v", err)
	}

	// Run migrations
	if err := DB.AutoMigrate(&models.User{}); err != nil {
		log.Fatalf("Migration failed: %v", err)
	}

	fmt.Println("Database is running!")
}
