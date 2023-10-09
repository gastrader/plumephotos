package controllers

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/gastrader/website/context"
	"github.com/gastrader/website/errors"
	"github.com/gastrader/website/models"
)

type Users struct {
	Templates struct {
		New    Template
		SignIn Template
		ForgotPassword Template
		CheckYourEmail Template
		ResetPassword Template
	}
	UserService    *models.UserService
	SessionService *models.SessionService
	PasswordResetService *models.PasswordResetService
	EmailService *models.EmailService
}

type UserData struct {
	Email    string
	Password string
}

func (u Users) New(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Email string
	}
	data.Email = r.FormValue("email")
	u.Templates.New.Execute(w, r, data)
}

func (u Users) Create(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Email string
		Password string
	}
	data.Email = r.FormValue("email")
	data.Password = r.FormValue("password")

	user, err := u.UserService.Create(data.Email, data.Password)
	if err != nil {
		if errors.Is(err, models.ErrEmailTaken){
			err = errors.Public(err, "That email address has been taken.")
		}
		u.Templates.New.Execute(w, r, data, err)
		return
	}
	session, err := u.SessionService.Create(user.ID)
	if err != nil {
		fmt.Println(err)
		http.Redirect(w, r, "/signin", http.StatusFound)
		return
	}
	fmt.Println("The session value is: ", session.Token)
	setCookie(w, CookieSession, session.Token)
	http.Redirect(w, r, "/galleries", http.StatusFound)
	fmt.Fprintf(w, "User Created: %+v", user)
}

func (u Users) SignIn(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Email string
	}
	data.Email = r.FormValue("email")
	u.Templates.SignIn.Execute(w, r, data)
}

func (u Users) ProcessSignIn(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Email    string
		Password string
	}
	data.Email = r.FormValue("email")
	data.Password = r.FormValue("password")
	user, err := u.UserService.Authenticate(data.Email, data.Password)
	if err != nil {
		if errors.Is(err, models.ErrLoginNotFound){
			err = errors.Public(err, "The email address and password do not match.")
			u.Templates.SignIn.Execute(w, r, data, err)
			return
		}
		if errors.Is(err, models.ErrEmailNotFound){
			err = errors.Public(err, "The email address was not found.")
			u.Templates.SignIn.Execute(w, r, data, err)
			return
		}
		fmt.Println(err)
		http.Error(w, "Login did not work .", http.StatusInternalServerError)

		return
	}
	session, err := u.SessionService.Create(user.ID)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Something went wrong...", http.StatusInternalServerError)
		return
	}
	setCookie(w, CookieSession, session.Token)
	http.Redirect(w, r, "/galleries", http.StatusFound)
	fmt.Fprintf(w, "User Signed in: %+v", user)
}

func (u Users) CurrentUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	user := context.User(ctx)
	fmt.Fprintf(w, "Current user: %s\n", user.Email)
}

func (u Users) ProcessSignOut(w http.ResponseWriter, r *http.Request) {
	token, err := readCookie(r, CookieSession)
	if err != nil {
		http.Redirect(w, r, "/signin", http.StatusFound)
		return
	}
	err = u.SessionService.Delete(token)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "something went wrong.", http.StatusInternalServerError)
		return
	}
	deleteCookie(w, CookieSession)
	http.Redirect(w, r, "/signin", http.StatusFound)
}

type UserMiddleware struct {
	SessionService *models.SessionService
}

//allows us to have URL like mysite.com/forgot-pw?email=jon@jon.com
func (u Users) ForgotPassword(w http.ResponseWriter, r *http.Request){
	var data struct {
		Email string
	}
	data.Email = r.FormValue("email")
	u.Templates.ForgotPassword.Execute(w, r, data)
}

func (u Users) ProcessForgotPassword(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Email string
	}
	data.Email = r.FormValue("email")
	pwReset, err := u.PasswordResetService.Create(data.Email)
	if err != nil{
		//handle other cases if user DNE w/ that email address
		fmt.Println(err)
		http.Error(w, "Something went wrong..", http.StatusInternalServerError)
		return
	}
	vals := url.Values{
		"token": {pwReset.Token},
	}
	resetURL := "https://www.plumephotos.com/reset-pw?"+ vals.Encode()
	err = u.EmailService.ForgotPassword(data.Email, resetURL)
	if err != nil{
		fmt.Println(err)
		http.Error(w, "Something went wrong..", http.StatusInternalServerError)
		return
	}
	//do not render reset token.
	u.Templates.CheckYourEmail.Execute(w, r, data)

}

func (u Users) ResetPassword(w http.ResponseWriter, r *http.Request){
	var data struct {
		Token string
	}
	data.Token = r.FormValue("token")
	u.Templates.ResetPassword.Execute(w, r, data)
}

func (u Users) ProcessResetPassword(w http.ResponseWriter, r *http.Request){
	var data struct {
		Token string
		Password string
	}
	data.Token = r.FormValue("token")
	data.Password = r.FormValue("password")
	
	user, err := u.PasswordResetService.Consume(data.Token)
	if err != nil{
		fmt.Println(err)
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}
	//Update user's password
	err = u.UserService.UpdatePassword(user.ID, data.Password)
	if err != nil{
		fmt.Println(err)
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}

	//Sign user in after pw reset. Same as "process sign in"
	session, err := u.SessionService.Create(user.ID)
	if err != nil{
		fmt.Println(err)
		http.Redirect(w, r, "/signin", http.StatusFound)
		return
	}
	setCookie(w, CookieSession, session.Token)
	http.Redirect(w, r, "/galleries", http.StatusFound)
}


func (umw UserMiddleware) SetUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, err := readCookie(r, CookieSession)
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}
		user, err := umw.SessionService.User(token)
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}
		ctx := r.Context()
		ctx = context.WithUser(ctx, user)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}

func (umw UserMiddleware) RequireUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
		user := context.User(r.Context())
		if user == nil {
			http.Redirect(w, r, "/signin", http.StatusFound)
			return
		}
		next.ServeHTTP(w, r)
	})
}

