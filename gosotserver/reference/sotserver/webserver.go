package main

import (
	"html/template"
	"net/http"
)

// for use with templates
type Page struct {
	Name  string
	Title string
}

var templates = template.Must(template.ParseGlob("templates/*.html"))

const viewsPath = "/"

func handle(t string) (string, http.HandlerFunc) {
	return viewsPath + t, func(w http.ResponseWriter, r *http.Request) {
		p := &Page{Name: "Leo Reyes", Title: t}
		err := templates.ExecuteTemplate(w, t+".html", p)
		if err != nil {
			http.NotFound(w, r)
		}
	}
}

func startWebServer() {
	http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("css/"))))
	http.Handle("/js/", http.StripPrefix("/js/", http.FileServer(http.Dir("js/"))))

	//http.HandleFunc(handle("index"))
	http.HandleFunc(handle("home"))
	http.HandleFunc(handle("devices"))
	http.HandleFunc(handle("homeAdmin"))
	http.HandleFunc(handle("mapAdmin"))
	http.HandleFunc(handle("mapUser"))
	http.HandleFunc(handle("users"))

	// our server is one line!
	http.ListenAndServe(":8080", nil)
}
