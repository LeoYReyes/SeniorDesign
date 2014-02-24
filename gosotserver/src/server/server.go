package main

import (
	"CustomRequest"
	"fmt"
	"webserver"
)

var toWebCh = make(chan *CustomRequest.Request)
var fromWebCh = make(chan *CustomRequest.Request)
var toDatabaseCh = make(chan *CustomRequest.Request)
var fromDatabaseCh = make(chan *CustomRequest.Request)
var toDeviceCh = make(chan *CustomRequest.Request)
var fromDeviceCh = make(chan *CustomRequest.Request)

func main() {
	// channel can take optional capacity param to make it asynchronous
	//comChannel := make(chan string)
	go webserver.StartWebServer(fromWebCh, toWebCh)
	fmt.Println(<-fromWebCh)
}
