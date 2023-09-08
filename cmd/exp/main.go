package main

import (
	"html/template"
	"os"
)

type User struct {
	Name string
	Age int
	Bio string
}

func main() {
	t, err := template.ParseFiles("hello.html")
	if err != nil {
		panic(err)
	}
	user := User{
		Name: "johnny",
		Age: 1232,
		Bio: `<script>alert("hacked")</script>`,
	}
	err = t.Execute(os.Stdout, user)
	if err != nil {
		panic(err)
	}

}
