package main

import (
	"html/template"
	"net/http"
	"os"
)

func serveDemo(w http.ResponseWriter, r *http.Request) {
	const fileName = "./templates/demo.html"
	t, _ := template.ParseFiles(fileName)
	data := struct{ APIUrl string }{os.Getenv("API_URL")}

	t.Execute(w, data)
}
