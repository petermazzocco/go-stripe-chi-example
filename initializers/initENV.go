package initializers

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func InitENV() {
	dir, _ := os.Getwd()
	log.Printf("Current working directory: %s", dir)

	// Try with explicit path
	err := godotenv.Load(dir + "/.env")
	if err != nil {
		log.Println("Warning: .env file not found or cannot be loaded - using environment variables instead")
	} else {
		log.Println(".env file loaded successfully")
	}
}
