package main

import (
	"fmt"

	"github.com/gastrader/website/models"
)


const (
	host     = "sandbox.smtp.mailtrap.io"
	port     = 25
	username = "297bd7f4a812a3"
	password = "b25328c4118822"
)

func main() {
	es := models.NewEmailService(models.SMTPConfig{
		Host:     host,
		Port:     port,
		Username: username,
		Password: password,
	})
	err := es.ForgotPassword("gavinp96@gmail.com", "https://localdankness.com")
	if err != nil {
		panic(err)

	}
	fmt.Println("email sent!")

}
