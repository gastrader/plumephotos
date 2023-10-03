package main

import (
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/gastrader/website/controllers"
	"github.com/gastrader/website/migrations"
	"github.com/gastrader/website/models"
	"github.com/gastrader/website/templates"
	"github.com/gastrader/website/views"
	"github.com/go-chi/chi/v5"
	"github.com/gorilla/csrf"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/joho/godotenv"
)

type config struct {
	PSQL models.PostgresConfig
	SMTP models.SMTPConfig
	CSRF struct {
		Key    string
		Secure bool
	}
	Server struct {
		Address string //":3000"
	}
}

func loadEnvConfig() (config, error) {
	var cfg config
	err := godotenv.Load()
	if err != nil {
		return cfg, err
	}
	cfg.PSQL = models.DefaultPostgresConfig()

	cfg.SMTP.Host = os.Getenv("SMTP_HOST")
	cfg.SMTP.Username = os.Getenv("SMTP_USERNAME")
	portStr := os.Getenv("SMTP_PORT")
	cfg.SMTP.Port, err = strconv.Atoi(portStr)
	if err != nil {
		return cfg, err
	}
	cfg.SMTP.Password = os.Getenv("SMTP_PASSWORD")

	cfg.CSRF.Key = "u2312casdyug682yubbcjyuihyu3bnsx"
	cfg.CSRF.Secure = false

	cfg.Server.Address = ":3000"
	return cfg, nil
}

func main() {
	cfg, err := loadEnvConfig()
	if err != nil {
		panic(err)
	}
	//setup the db

	db, err := models.Open(cfg.PSQL)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	err = models.MigrateFS(db, migrations.FS, ".")
	if err != nil {
		panic(err)
	}

	//setup services
	userService := &models.UserService{
		DB: db,
	}
	sessionService := &models.SessionService{
		DB: db,
	}
	pwResetService := &models.PasswordResetService{
		DB: db,
	}
	emailService := models.NewEmailService(cfg.SMTP)

	//setup middleware
	umw := controllers.UserMiddleware{
		SessionService: sessionService,
	}

	csrfMw := csrf.Protect([]byte(cfg.CSRF.Key), csrf.Secure(cfg.CSRF.Secure))

	//setup controllers
	usersC := controllers.Users{
		UserService:          userService, //TODO = set this
		SessionService:       sessionService,
		PasswordResetService: pwResetService,
		EmailService:         emailService,
	}
	usersC.Templates.New = views.Must(views.ParseFS(templates.FS, "signup.html", "tailwind.html"))
	usersC.Templates.SignIn = views.Must(views.ParseFS(templates.FS, "signin.html", "tailwind.html"))
	usersC.Templates.ForgotPassword = views.Must(views.ParseFS(templates.FS, "forgot_pw.html", "tailwind.html"))
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
	r.Get("/forgot-pw", usersC.ForgotPassword)
	r.Post("/forgot-pw", usersC.ProcessForgotPassword)
	r.Get("/users/me", usersC.CurrentUser)
	r.Route("/users/me", func(r chi.Router) {
		r.Use(umw.RequireUser)
		r.Get("/", usersC.CurrentUser)
	})
	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Page Not Found", http.StatusNotFound)
	})

	fmt.Printf("starting the server on %s..\n", cfg.Server.Address)
	err = http.ListenAndServe(cfg.Server.Address, r)
	if err != nil{
		panic(err)
	}

}
