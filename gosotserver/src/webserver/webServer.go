package webserver

import (
	"CustomProtocol"
	"crypto/sha1"
	//"databaseSOT"
	"fmt"
	"html/template"
	"mux"
	"net/http"
	"strings"
)

// for use with templates
type Page struct {
	UserName      string
	DeviceName    string
	DeviceStatus  string
	KeyloggerText string
}

var templates = template.Must(template.ParseGlob("templates/*.html"))

const viewsPath = "/"

func handle(t string) (string, http.HandlerFunc) {
	return viewsPath + t, func(w http.ResponseWriter, r *http.Request) {
		//fmt.Println("url: ", t)
		// check for session. If no session is active, redirect to home page
		_, cookieErr := r.Cookie("userSession")
		sesh, _ := store.Get(r, "userSession")
		if cookieErr != nil {
			//fmt.Println("Cookie Error: ", cookieErr)
		} else {
			//fmt.Println("Cookie userid: ", sesh.Values["userId"])
			//fmt.Println("Cookie isLoggedIn: ", sesh.Values["isLoggedIn"])
		}
		if t != "home" {
			if sesh.Values["isLoggedIn"] == "true" {
				userName := string(sesh.Values["userId"].(string))
				p := &Page{UserName: userName, DeviceName: "Test Device", DeviceStatus: "Not Stolen", KeyloggerText: "Test Key Log"}
				err := templates.ExecuteTemplate(w, t+".html", p)
				if err != nil {
					fmt.Println(err)
					//http.NotFound(w, r)
				}
			} else {
				http.Redirect(w, r, "home", http.StatusFound)
			}
		} else {
			if sesh.Values["isLoggedIn"] == "true" {
				http.Redirect(w, r, "/UserMapNEW", http.StatusFound)
			}
			p := &Page{}
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
	loginName := strings.Trim(r.PostForm.Get("loginName"), " ")
	loginPassword := strings.Trim(r.PostForm.Get("loginPassword"), " ")

	if loginName == "" || loginPassword == "" {
		// Send error message back
		return
	}
	//TODO: move hashing to database VerifyAccountFunction
	ha := sha1.New()
	ha.Write([]byte(strings.Join([]string{loginName, loginPassword}, "")))

	/*fmt.Println(loginName)
	fmt.Println(loginPassword)*/
	hashedPass := fmt.Sprintf("%x", ha.Sum(nil))

	buf := []byte{}
	buf = append(buf, []byte(loginName)...)
	buf = append(buf, 0x1B)
	buf = append(buf, []byte(hashedPass)...)
	buf = append(buf, 0x1B)

	resCh := make(chan []byte)
	req := &CustomProtocol.Request{Id: CustomProtocol.AssignRequestId(), Destination: CustomProtocol.Database, Source: CustomProtocol.Web,
		OpCode: CustomProtocol.VerifyLoginCredentials, Payload: buf, Response: resCh}
	toServer <- req
	sucessful, res := CustomProtocol.GetResponse(resCh, 10)
	//res := <-resCh
	if sucessful {
		//fmt.Println("Response: ", res[0], res[1])
		if (res[0] == 1) && (res[1] == 1) {
			serveSession(w, r)

		} else {
			http.Error(w, "Invalid Login", 80085)
			return
		}
	} else {
		//fmt.Println("No response on login")
		http.Error(w, "Timeout on Login", 8008135)
		return
	}
}

func logoutHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "userSession")
	session.Values["isLoggedIn"] = "false"
	session.Options.MaxAge = -1
	session.Save(r, w)
	http.Redirect(w, r, "/home", http.StatusFound)
}

func signUpHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		// Handle error
		fmt.Println("ParseForm error: ", err)
	}
	//TODO: check for blank inputs
	firstName := strings.Trim(r.PostForm.Get("firstName"), " ")
	lastName := strings.Trim(r.PostForm.Get("lastName"), " ")
	loginName := strings.Trim(r.PostForm.Get("loginName"), " ")
	phoneNumber := strings.Trim(r.PostForm.Get("phoneNumber"), " ")
	password := strings.Trim(r.PostForm.Get("password"), " ")

	if firstName == "" || lastName == "" || loginName == "" || phoneNumber == "" || password == "" {
		// Redirect back to home if any fields are empty, someone is trying to play tricks
		http.Redirect(w, r, "/home", http.StatusFound)
		return
	}
	//TODO: move hashing to database SignUp function
	h := sha1.New()
	h.Write([]byte(strings.Join([]string{loginName, password}, "")))

	/*fmt.Println(firstName)
	fmt.Println(lastName)
	fmt.Println(loginName)
	fmt.Println(phoneNumber)
	fmt.Println(password)*/
	hashedPass := fmt.Sprintf("%x", h.Sum(nil))
	//fmt.Println(hashedPass)

	buf := []byte{}
	buf = append(buf, []byte(firstName)...)
	buf = append(buf, 0x1B)
	buf = append(buf, []byte(lastName)...)
	buf = append(buf, 0x1B)
	buf = append(buf, []byte(loginName)...)
	buf = append(buf, 0x1B)
	buf = append(buf, []byte(phoneNumber)...)
	buf = append(buf, 0x1B)
	buf = append(buf, []byte(hashedPass)...)
	buf = append(buf, 0x1B)

	resCh := make(chan []byte)
	req := &CustomProtocol.Request{Id: CustomProtocol.AssignRequestId(), Destination: CustomProtocol.Database, Source: CustomProtocol.Web,
		OpCode: CustomProtocol.NewAccount, Payload: buf, Response: resCh}
	toServer <- req
	successful, res := CustomProtocol.GetResponse(resCh, 10)
	if successful {
		//fmt.Println("Response: ", res[0])
		//databaseSOT.SignUp(r.PostForm.Get("firstName"), r.PostForm.Get("lastName"),
		//	r.PostForm.Get("email"), r.PostForm.Get("phoneNumber"), hashedPass)
		if res[0] == 1 {
			serveSession(w, r)
		} else {
			http.Error(w, "Duplicate Account", 7000)
		}
	} else {
		http.Error(w, "Timeout on Signup", 7001)
	}
}

// Declaration of global variable
var toServer chan *CustomProtocol.Request
var fromServer chan *CustomProtocol.Request

func StartWebServer(toServerIn chan *CustomProtocol.Request, fromServerIn chan *CustomProtocol.Request) {
	http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("css/"))))
	http.Handle("/js/", http.StripPrefix("/js/", http.FileServer(http.Dir("js/"))))
	http.Handle("/images/", http.StripPrefix("/images/", http.FileServer(http.Dir("images/"))))
	http.Handle("/media/", http.StripPrefix("/media/", http.FileServer(http.Dir("media/"))))

	toServer = toServerIn
	fromServer = fromServerIn

	go h.run()
	go chanHandler()

	r := mux.NewRouter()

	r.HandleFunc(handle(""))
	r.HandleFunc(handle("home"))
	r.HandleFunc(handle("mapUser"))
	r.HandleFunc(handle("homeAdmin"))
	r.HandleFunc(handle("mapAdmin"))
	r.HandleFunc(handle("webSocketTest"))
	r.HandleFunc(handle("messager"))
	r.HandleFunc(handle("devices"))
	r.HandleFunc(handle("users"))
	r.HandleFunc(handle("UserMapNEW"))
	r.HandleFunc("/signup", signUpHandler)
	r.HandleFunc("/login", loginHandler)
	r.HandleFunc("/logout", logoutHandler)
	r.HandleFunc("/ws", serveWs)
	r.HandleFunc("/newDevice", newDeviceHandler)
	r.HandleFunc("/getDeviceInfo", deviceInfoHandler)
	r.HandleFunc("/toggleDevice", toggleDeviceHandler)
	r.HandleFunc("/pingDevice", pingDeviceHandler)

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
			fmt.Println("web server received: ", req.Payload)
			//TODO: parse payload and send coordinates to correct ws session
			//TODO: create process request
			//parsedPayload := CustomProtocol.ParsePayload(req.Payload)

			/*msg := []byte{}
			// DeviceId, Phone Number
			msg = append(msg, []byte(parsedPayload[0])...)
			msg = append(msg, 0x1B)
			// Latitude
			msg = append(msg, []byte(parsedPayload[1])...)
			msg = append(msg, 0x1B)
			// Longitude
			msg = append(msg, []byte(parsedPayload[2])...)
			msg = append(msg, 0x1B)

			fmt.Println(msg)*/

			h.broadcast <- req.Payload
		}
	}
}
