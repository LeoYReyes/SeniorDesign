package main

import (
	"webserver"
)

func main() {
	go webserver.StartWebServer()
}
