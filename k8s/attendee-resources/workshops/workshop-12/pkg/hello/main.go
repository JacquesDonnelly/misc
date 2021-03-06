package hello

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"time"
)

type greeter struct {
	version string
	greetee string
}

func init() {
}

func (g *greeter) rootHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, %s from %s!\n", g.greetee, r.RemoteAddr)
	fmt.Fprintf(w, "Welcome to %s\n", r.URL.Path)
	fmt.Fprintf(w, "Version %s\n", g.version)
}

func serveAdmin() {
	adminRouter := http.NewServeMux()

	adminRouter.HandleFunc("/healthz", healthzHandler)
	adminRouter.HandleFunc("/ready", readyHandler)
	adminRouter.HandleFunc("/set-unhealthy", setUnhealthyHandler)
	adminRouter.HandleFunc("/set-healthy", setHealthyHandler)
	adminRouter.HandleFunc("/set-unready", setUnreadyHandler)
	adminRouter.HandleFunc("/set-ready", setReadyHandler)

	adminRouter.HandleFunc("/pre-stop", preStopHandler)
	adminRouter.HandleFunc("/stopping-state", stoppingStateHandler)

	adminServer := http.Server{
		Addr:    ":8081",
		Handler: adminRouter,
	}
	fmt.Println("Serving health endpoints on :8081")
	adminServer.ListenAndServe()
}

func RunServer(version string, name string) {
	g := &greeter{
		version: version,
		greetee: name,
	}

	go serveAdmin()
	go delayReady()

	helloRouter := http.NewServeMux()
	helloRouter.HandleFunc("/", g.rootHandler)
	helloServer := http.Server{
		Addr:    ":8080",
		Handler: helloRouter,
	}
	fmt.Println("Serving on :8080")
	helloServer.ListenAndServe()
}

func RunTask(endpoint string) {
	// perform 'work'
	fmt.Println("Sleeping for 5 seconds, imagine I'm working hard on a business critical job...")
	time.Sleep(5 * time.Second)

	message := time.Now().UTC().Format(time.RFC850)
	_, err := http.PostForm(endpoint, url.Values{"message": {"hello"}, "time": {message}})
	if err != nil {
		log.Fatal(err)
	}
}
