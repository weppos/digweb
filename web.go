package main

import (
	"fmt"
	"net/http"
	"os"
)

var (
	serverPort = os.Getenv("PORT")
)

func main() {
	http.HandleFunc("/", RootHandler)

	fmt.Println(fmt.Sprintf("Listening on %s...", serverPort))
	err := http.ListenAndServe(":"+serverPort, nil)
	if err != nil {
		panic(err)
	}
}

// RootHandler handles HTTP requests to /.
func RootHandler(res http.ResponseWriter, req *http.Request) {

	// The handler for / effectively becomes a catch-all
	// for all the not-handled routes.
	// I don't want the app to respond to any possible PATH,
	// hence let's kill invalid requests immediately.
	if req.URL.Path != "/" {
		CatchallHandler(res, req)
		return
	}

	switch req.Method {
	case "GET":
		actionRoot(res, req)
	case "POST":
		actionDig(res, req)
	default:
		http.NotFound(res, req)
	}
}

// CatchallHandler handles all unhandled routes.
func CatchallHandler(res http.ResponseWriter, req *http.Request) {
	http.NotFound(res, req)
}

// actionRoot responds with a simple Alive message.
func actionRoot(res http.ResponseWriter, req *http.Request) {
	fmt.Fprintln(res, "Alive!")
}

// actionDig executes the dig query and responds with the result.
func actionDig(res http.ResponseWriter, req *http.Request) {
	fmt.Fprintln(res, "Digging...")
}
