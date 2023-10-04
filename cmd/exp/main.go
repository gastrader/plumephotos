package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/gastrader/website/models"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env")
	}
	host := os.Getenv("SMTP_HOST")
	portStr := os.Getenv("SMTP_PORT")
	port, err := strconv.Atoi(portStr)
	if err != nil {
		panic(err)
	}
	username := os.Getenv("SMTP_USERNAME")
	password := os.Getenv("SMTP_PASSWORD")
	es := models.NewEmailService(models.SMTPConfig{
		Host:     host,
		Port:     port,
		Username: username,
		Password: password,
	})
	err = es.ForgotPassword("gavinp96@gmail.com", "https://localhost.com")
	if err != nil {
		panic(err)

	}
	fmt.Println("email sent!")
	// gs := models.GalleryService{}
	// fmt.Println(gs.Images(7))
}
