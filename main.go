package main

import (
	"fmt"
	"net/http"

	"github.com/gastrader/website/controllers"
	"github.com/gastrader/website/models"
	"github.com/gastrader/website/templates"
	"github.com/gastrader/website/views"
	"github.com/go-chi/chi/v5"
	"github.com/gorilla/csrf"
	_ "github.com/jackc/pgx/v4/stdlib"
)

func main() {
	r := chi.NewRouter()

	tpl := views.Must(views.ParseFS(templates.FS, "home.html", "tailwind.html"))
	r.Get("/", controllers.StaticHandler(tpl))

	tpl = views.Must(views.ParseFS(templates.FS, "contact.html", "tailwind.html"))
	r.Get("/contact", controllers.FAQ(tpl))

	cfg := models.DefaultPostgresConfig()
	db, err := models.Open(cfg)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	userService := models.UserService{
		DB: db,
	}
	usersC := controllers.Users{
		UserService: &userService, //TODO = set this
	}
	//TEMPLATE PARSING
	usersC.Templates.New = views.Must(views.ParseFS(
		templates.FS,
		"signup.html", "tailwind.html"))
	usersC.Templates.SignIn = views.Must(views.ParseFS(
		templates.FS,
		"signin.html", "tailwind.html"))
		
	//ROUTE HANDLING
	r.Get("/signup", usersC.New)
	r.Post("/users", usersC.Create)
	r.Get("/signin", usersC.SignIn)
	r.Post("/signin", usersC.ProcessSignIn)
	r.Get("/user/me", usersC.CurrentUser)

	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Page notttt Found", http.StatusNotFound)
	})

	fmt.Println("starting the server on :3000..")

	csrfKey := "u2312casdyug682yubbcjyuihyu3bnsx"
	csrfMw := csrf.Protect([]byte(csrfKey), csrf.Secure(false))
	
	http.ListenAndServe("127.0.0.1:3000", csrfMw(r))

}
