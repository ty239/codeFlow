package main

import ( // lets us create packages and import them to are codeBase
	"fmt" // package that used to create http server send request etc
	"net/http"
)

func homeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "hello backend") //w send a response/data to the server
	// while r contains the information about the request
}
