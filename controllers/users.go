package controllers

import (
	"fmt"
	"net/http"

	"github.com/gastrader/website/models"
)

type Users struct {
	Templates struct {
		New Template
	}
	UserService *models.UserService
}

type UserData struct{
	Email string
	Password string
}

func (u Users) New(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Email string
	}
	data.Email = r.FormValue("email")
	u.Templates.New.Execute(w, data)
}

//COULD USE BODYPARSER FROM GORILLA - PART OF FIBER...
func (u Users) Create(w http.ResponseWriter, r *http.Request){
	email := r.FormValue("email")
	password := r.FormValue("password")
	
	user, err := u.UserService.Create(email, password)
	if err != nil{
		fmt.Println(err)
		http.Error(w, "Login did not work .", http.StatusInternalServerError)
		return 
	}
	fmt.Fprintf(w, "User Created: %+v", user)
}