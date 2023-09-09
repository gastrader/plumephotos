package main

import (
	"fmt"
	"net/http"

	"github.com/gastrader/website/controllers"
	"github.com/gastrader/website/templates"
	"github.com/gastrader/website/views"
	"github.com/go-chi/chi/v5"
)

func main() {
	r := chi.NewRouter()

	tpl := views.Must(views.ParseFS(templates.FS, "home.html", "tailwind.html")) 
	r.Get("/", controllers.StaticHandler(tpl))

	tpl = views.Must(views.ParseFS(templates.FS, "contact.html", "tailwind.html"))
	r.Get("/contact", controllers.FAQ(tpl))

	//MUST is panicing if there is error that template being rendered.
	usersC := controllers.Users{}
	usersC.Templates.New = views.Must(views.ParseFS(
		templates.FS,
		"signup.html", "tailwind.html"))
	r.Get("/signup", usersC.New)
	r.Post("/users", usersC.Create)

	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Page notttt Found", http.StatusNotFound)
	})

	fmt.Println("starting the server on :3333...")
	http.ListenAndServe("127.0.0.1:3333", r)
}
