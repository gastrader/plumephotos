package controllers

import (
	"fmt"
	"net/http"
)

type Users struct {
	Templates struct {
		New Template
	}
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
	fmt.Fprintln(w, r.FormValue("email"))
	fmt.Fprint(w, r.FormValue("password"))
}