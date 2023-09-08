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

	tpl := views.Must(views.ParseFS(templates.FS, "home.html"))
	r.Get("/", controllers.StaticHandler(tpl))

	tpl = views.Must(views.ParseFS(templates.FS, "contact.html"))
	r.Get("/contact", controllers.StaticHandler(tpl))

	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Page notttt Found", http.StatusNotFound)
	})

	fmt.Println("starting the server on :3333...")
	http.ListenAndServe("127.0.0.1:3333", r)
}
