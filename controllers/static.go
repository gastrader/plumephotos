package controllers

import (
	"html/template"
	"net/http"

	"github.com/gastrader/website/views"
)

func StaticHandler(tpl views.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tpl.Execute(w, nil)
	}
}

func FAQ(tpl views.Template) http.HandlerFunc {
	questions := []struct {
		Question string
		Answer   template.HTML
	}{
		{
			Question: "Is it free?",
			Answer:   "No its not",
		},
		{
			Question: "Is it not?",
			Answer:   "No its def not!",
		},
		{
			Question: "Where da office?",
			Answer:   "Remote mny guy!",
		},
		{
			Question: "How do I contact you bro?",
			Answer:   `Please do not contact me at <a href="mailto:gaslimits@gmail.com">gaslimits@gmail.com</a>`,
		},
	}
	return func(w http.ResponseWriter, r *http.Request) {
		tpl.Execute(w, questions)
	}

}
