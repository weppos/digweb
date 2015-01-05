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
	http.HandleFunc("/slack", SlackHandler)

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

// SlackHandler handles HTTP requests to /slack.
func SlackHandler(res http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case "POST":
		actionSlack(res, req)
	default:
		http.NotFound(res, req)
	}
}

// CatchallHandler handles all unhandled HTTP requests.
func CatchallHandler(res http.ResponseWriter, req *http.Request) {
	http.NotFound(res, req)
}

// actionRoot responds with a simple Alive message.
func actionRoot(res http.ResponseWriter, req *http.Request) {
	fmt.Fprintln(res, "Alive!")
}

// actionDig executes a dig query extracting the args from the request body
// and responds with the result.
func actionDig(res http.ResponseWriter, req *http.Request) {
	args, _ := ioutil.ReadAll(req.Body)
	writeDig(res, string(args))
}

// actionSlack executes a dig query extracting the args from a Slack payload
// and responds with the result.
func actionSlack(res http.ResponseWriter, req *http.Request) {
	args := req.FormValue("text")
	writeDig(res, string(args))
}

// Dig [@global-server] [domain] [q-type] [q-class] {q-opt}
func Dig(args string) ([]byte, error) {
	return exec.Command("dig", args).CombinedOutput()
}

func writeDig(res http.ResponseWriter, args string) {
	out, err := Dig(args)

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
