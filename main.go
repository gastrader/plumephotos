package main

import (
	"fmt"
	"net/http"

	"github.com/gastrader/website/controllers"
	"github.com/gastrader/website/migrations"
	"github.com/gastrader/website/models"
	"github.com/gastrader/website/templates"
	"github.com/gastrader/website/views"
	"github.com/go-chi/chi/v5"
	"github.com/gorilla/csrf"
	_ "github.com/jackc/pgx/v4/stdlib"
)

func main() {
	//setup the db
	cfg := models.DefaultPostgresConfig()
	db, err := models.Open(cfg)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	err = models.MigrateFS(db, migrations.FS, ".")
	if err != nil {
		panic(err)
	}

	//setup services
	userService := models.UserService{
		DB: db,
	}
	sessionService := models.SessionService{
		DB: db,
	}

	//setup middleware
	umw := controllers.UserMiddleware{
		SessionService: &sessionService,
	}
	csrfKey := "u2312casdyug682yubbcjyuihyu3bnsx"
	csrfMw := csrf.Protect([]byte(csrfKey), csrf.Secure(false))

	//setup controllers
	usersC := controllers.Users{
		UserService:    &userService, //TODO = set this
		SessionService: &sessionService,
	}
	usersC.Templates.New = views.Must(views.ParseFS(
		templates.FS,
		"signup.html", "tailwind.html"))
	usersC.Templates.SignIn = views.Must(views.ParseFS(
		templates.FS,
		"signin.html", "tailwind.html"))

	//setup router and routes
	r := chi.NewRouter()
	r.Use(csrfMw)
	r.Use(umw.SetUser)
	tpl := views.Must(views.ParseFS(templates.FS, "home.html", "tailwind.html"))
	r.Get("/", controllers.StaticHandler(tpl))

	tpl = views.Must(views.ParseFS(templates.FS, "contact.html", "tailwind.html"))
	r.Get("/contact", controllers.FAQ(tpl))

	//ROUTE HANDLING
	r.Get("/signup", usersC.New)
	r.Post("/users", usersC.Create)
	r.Get("/signin", usersC.SignIn)
	r.Post("/signin", usersC.ProcessSignIn)
	r.Post("/signout", usersC.ProcessSignOut)
	r.Get("/users/me", usersC.CurrentUser)
	r.Route("/users/me", func(r chi.Router){
		r.Use(umw.RequireUser)
		r.Get("/", usersC.CurrentUser)
	})
	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Page Not Found", http.StatusNotFound)
	})

	fmt.Println("starting the server on :3000..")
	http.ListenAndServe("127.0.0.1:3000", r)

}
