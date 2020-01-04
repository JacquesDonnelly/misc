package hello

import (
	"fmt"
	"net/http"
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

// Run starts the HTTP server
func Run(version string, name string) {
	g := &greeter{
		version: version,
		greetee: name,
	}

	helloRouter := http.NewServeMux()
	helloRouter.HandleFunc("/", g.rootHandler)
	helloServer := http.Server{
		Addr:    ":8080",
		Handler: helloRouter,
	}
	fmt.Println("Serving on :8080")
	helloServer.ListenAndServe()
}
