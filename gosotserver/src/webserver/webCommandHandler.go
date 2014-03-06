package webserver

import (
	"fmt"
	"net/http"
)

func newDeviceHandler(w http.ResponseWriter, r *http.Request) {
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
