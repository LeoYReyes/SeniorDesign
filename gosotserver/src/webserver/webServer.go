package webserver

import (
	"CustomProtocol"
	"crypto/sha1"
	"databaseSOT"
	"fmt"
	"html/template"
	"mux"
	"net/http"
	"strconv"
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
		fmt.Println("url: ", t)
		// check for session. If no session is active, redirect to home page
		_, cookieErr := r.Cookie("userSession")
		sesh, _ := store.Get(r, "userSession")
		if cookieErr != nil {
			fmt.Println("Cookie Error: ", cookieErr)
		} else {
			fmt.Println("Cookie userid: ", sesh.Values["userId"])
			fmt.Println("Cookie isLoggedIn: ", sesh.Values["isLoggedIn"])
		}
		p := &Page{Name: "Leo Reyes", Title: t}
		if t != "home" {
			if sesh.Values["isLoggedIn"] == "true" {
				err := templates.ExecuteTemplate(w, t+".html", p)
				if err != nil {
					http.NotFound(w, r)
				}
			} else {
				http.Error(w, "Not Logged In", 000)
			}
		} else {
			err := templates.ExecuteTemplate(w, "home.html", p)
			if err != nil {
				http.NotFound(w, r)
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

	fmt.Println(r.PostForm)
	if err != nil {
		// Handle error
		fmt.Println("decoder error: ", err)
	}

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

	buf := []byte{}
	buf = strconv.AppendQuoteToASCII(buf, r.PostForm.Get("loginName"))
	buf = append(buf, 0x1B)
	buf = strconv.AppendQuoteToASCII(buf, hashedPass)
	buf = append(buf, 0x1B)

	resCh := make(chan []byte)
	req := &CustomProtocol.Request{Id: CustomProtocol.AssignRequestId(), Destination: 1, Source: 2,
		OpCode: CustomProtocol.VerifyLoginCredentials, Payload: buf, Response: resCh}
	toServer <- req
	res := <-resCh
	//accountValid, passwordValid := databaseSOT.VerifyAccountInfo(r.PostForm.Get("loginName"), hashedPass)
	fmt.Println("Response: ", res[0], res[1])
	if (res[0] == 1) && (res[1] == 1) {
		serveSession(w, r)
	} else {
		http.Error(w, "Invalid Login", 80085)
		return
	}
}

func logoutHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "userSession")
	session.Values["isLoggedIn"] = "false"
	session.Save(r, w)
	http.Redirect(w, r, "/home", http.StatusFound)
}

func signUpHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		// Handle error
		fmt.Println("ParseForm error: ", err)
	}
	h := sha1.New()
	h.Write([]byte(strings.Join([]string{r.PostForm.Get("email"), r.PostForm.Get("password")}, "")))

	fmt.Println(r.PostForm.Get("firstName"))
	fmt.Println(r.PostForm.Get("lastName"))
	fmt.Println(r.PostForm.Get("email"))
	fmt.Println(r.PostForm.Get("phoneNumber"))
	fmt.Println(r.PostForm.Get("password"))
	hashedPass := fmt.Sprintf("%x", h.Sum(nil))
	fmt.Println(hashedPass)
	databaseSOT.SignUp(r.PostForm.Get("firstName"), r.PostForm.Get("lastName"),
		r.PostForm.Get("email"), r.PostForm.Get("phoneNumber"), hashedPass)
	serveSession(w, r)
}

// Declaration of global variable
var toServer chan *CustomProtocol.Request
var fromServer chan *CustomProtocol.Request

func StartWebServer(toServerIn chan *CustomProtocol.Request, fromServerIn chan *CustomProtocol.Request) {
	http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("css/"))))
	http.Handle("/js/", http.StripPrefix("/js/", http.FileServer(http.Dir("js/"))))
	http.Handle("/images/", http.StripPrefix("/images/", http.FileServer(http.Dir("images/"))))

	toServer = toServerIn
	fromServer = fromServerIn

	go h.run()
	go chanHandler()

	r := mux.NewRouter()

	r.HandleFunc(handle("home"))
	r.HandleFunc(handle("mapUser"))
	r.HandleFunc(handle("homeAdmin"))
	r.HandleFunc(handle("mapAdmin"))
	r.HandleFunc(handle("webSocketTest"))
	r.HandleFunc(handle("messager"))
	r.HandleFunc(handle("devices"))
	r.HandleFunc(handle("users"))
	r.HandleFunc("/signup", signUpHandler)
	r.HandleFunc("/login", loginHandler)
	r.HandleFunc("/logout", logoutHandler)
	r.HandleFunc("/ws", serveWs)

	http.Handle("/", r)
	// our server is one line!
	http.ListenAndServe(":8080", nil)

}

//TODO: make chan handler in webClientHub
func chanHandler() {
	fmt.Println("web chan handler started")
	for {
		select {
		case req := <-fromServer:
			fmt.Println("web server received: ", string(req.Payload))
			h.broadcast <- req.Payload
		}
	}
}
