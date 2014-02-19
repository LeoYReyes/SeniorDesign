package webserver

import (
	"fmt"
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

func Test() {
	fmt.Println("test")
}

func StartWebServer() {
	http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("css/"))))
	http.Handle("/js/", http.StripPrefix("/js/", http.FileServer(http.Dir("js/"))))

	go h.run()

	http.HandleFunc(handle("mapUser"))
	http.HandleFunc(handle("webSocketTest"))
	http.HandleFunc(handle("messager"))
	http.HandleFunc("/ws", serveWs)

	// our server is one line!
	http.ListenAndServe(":8080", nil)
	if err := http.ListenAndServe(":1234", nil); err != nil {
		fmt.Println("ListenAndServer: ", err)
	}
}
