package main

import (
	"webserver"
)

func main() {
	// channel can take optional capacity param to make it asynchronous
	//comChannel := make(chan string)
	toWebSocketCh := make(chan Request)
	fromWebSocketCh := make(chan Request)
	toDatabaseCh := make(chan Request)
	fromDatabaseCh := make(chan Request)
	toDeviceCh := make(chan Request)
	fromDeviceCh := make(chan Request)

	webserver.StartWebServer()
}
