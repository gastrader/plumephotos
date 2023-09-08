package main

import (
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/contact", contactHandler)
	fmt.Println("starting the server on :3000...")
	http.ListenAndServe(":3000", nil)
}

func homeHandler(w http.ResponseWriter, r *http.Request){
	w.Header().Add("Content-Type", "text/html")
	fmt.Fprint(w, "<h1>Hello there matey!</h1>")
}

func contactHandler(w http.ResponseWriter, r *http.Request){
	w.Header().Add("Content-Type", "text/html")
	fmt.Fprint(w, "<h1>CONTACT PAGE</h1><p>To get in touch, email me at: <a href=\"mailto:gaslimits@gmail.com\">gavinp96@gmail.com</a></p>")
}