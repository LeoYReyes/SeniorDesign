package main

import (
	"webserver"
)

func main() {
	comChannel := make(chan string)
	go webserver.StartWebServer()
}
