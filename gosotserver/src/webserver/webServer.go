package webserver

import (
	"crypto/sha1"
	"databaseSOT"
	"fmt"
	"html/template"
	"mux"
	"net/http"
	"strings"
	//"schema"
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
		fmt.Println(t)
		fmt.Println(r)
		p := &Page{Name: "Leo Reyes", Title: t}
		if t == "home" {
			err := templates.ExecuteTemplate(w, t+".html", p)
			if err != nil {
				http.NotFound(w, r)
			}
		} else {
			if r.Header.Get("Origin") != "http://"+r.Host {
				http.Error(w, "Not logged in", 8008)
			}
		}
	}
}

func formHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		// Handle error
		fmt.Println("ParseForm error: ", err)
	}
	//workorder := new(WorkOrder)
	// r.PostForm is a map of our POST form values
	//err = decoder.Decode(workorder, r.PostForm)
	fmt.Println(r.PostForm)
	//workorder.contactInfo.firstName = r.PostForm.Get("contactInfo.firstName")
	if err != nil {
		// Handle error
		fmt.Println("decoder error: ", err)
	}
	//fmt.Println(workorder)
	// Do something workorder
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		// Handle error
		fmt.Println("ParseForm error: ", err)
	}
	h := sha1.New()
	h.Write([]byte(strings.Join([]string{r.PostForm.Get("loginName"), r.PostForm.Get("loginPassword")}, "")))

	fmt.Println(r.PostForm.Get("loginName"))
	fmt.Println(r.PostForm.Get("loginPassword"))
	hashedPass := fmt.Sprintf("%x", h.Sum(nil))
	accountValid, passwordValid := databaseSOT.VerifyAccountInfo(r.PostForm.Get("loginName"), hashedPass)
	if accountValid && passwordValid {
		serveSession(w, r)
		http.Redirect(w, r, "/home/"+r.PostForm.Get("loginName")+"/mapUser", 301)
	} else {
		http.Error(w, "Invalid Login", 80085)
		return
	}
}

func StartWebServer() {
	http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("css/"))))
	http.Handle("/js/", http.StripPrefix("/js/", http.FileServer(http.Dir("js/"))))

	go h.run()

	r := mux.NewRouter()
	//vars := mux.Vars(request)
	//userId := vars["userid"]
	//s := r.PathPrefix("/home/{userid}").Subrouter()
	//s.HandleFunc(handle("/")).Name("userId")
	//url, err := s.Get("userId").URL("userid")

	r.HandleFunc(handle("home"))
	r.HandleFunc(handle("mapUser"))
	r.HandleFunc(handle("homeAdmin"))
	r.HandleFunc(handle("mapAdmin"))
	r.HandleFunc(handle("webSocketTest"))
	r.HandleFunc(handle("messager"))
	r.HandleFunc(handle("devices"))
	r.HandleFunc(handle("users"))
	r.HandleFunc("/login", loginHandler)
	r.HandleFunc("/ws", serveWs)

	http.Handle("/", r)

	// our server is one line!
	http.ListenAndServe(":8080", nil)

}
