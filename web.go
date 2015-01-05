package main

import (
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"io/ioutil"
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
	args, _ := ioutil.ReadAll(req.Body)
	out, err := exec.Command("dig", string(args)).CombinedOutput()

	if err != nil {
		http.Error(res, "Bad Request", http.StatusBadRequest)
		fmt.Fprintln(res, err)
	} else {
		// For now, I don't care whether the output is Stdout or Stderr.
		// In the future, we may want more control over the request and return
		// - 200 if Stdout
		// - 520 (Origin Error) if Stderr
		fmt.Fprintln(res, string(out))
	}
}
