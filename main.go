package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	startServer()
}

func startServer() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", homeHandler) //when someone vists / the function homeHandler will extucte

	fmt.Println("server is gonna run on port 8080") // health check of the server
	log.Fatal(http.ListenAndServe(":8080", mux))     // listen for requests and do it on the mux router
}
