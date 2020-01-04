package hello

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/go-redis/redis"
)

type greeter struct {
	version string
	greetee string
}

var redisClient *redis.Client

func init() {
}

func (g *greeter) rootHandler(w http.ResponseWriter, r *http.Request) {
	count, err := redisClient.Incr("view-counter").Result()
	var countMessage string
	if err != nil {
		countMessage = fmt.Sprintf("Unable to get view count (%s)", err)
	} else {
		countMessage = fmt.Sprintf("Views: %d", count)
	}

	fmt.Fprintf(w, "Hello, %s from %s!\n", g.greetee, r.RemoteAddr)
	fmt.Fprintf(w, "Welcome to %s\n", r.URL.Path)
	fmt.Fprintf(w, "Version %s\n", g.version)
	fmt.Fprintln(w, countMessage)
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

func RunServer(version string, name string, redisUrl string) {
	g := &greeter{
		version: version,
		greetee: name,
	}

	redisClient = redis.NewClient(&redis.Options{
		Addr:     redisUrl,
		Password: "",
		DB:       0,
	})

	go serveAdmin()
	go delayReady()

	helloRouter := http.NewServeMux()
	helloRouter.HandleFunc("/favicon.ico", http.NotFound)
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
	time.Sleep(20 * time.Second)

	message := time.Now().UTC().Format(time.RFC850)
	_, err := http.PostForm(endpoint, url.Values{"message": {"hello"}, "time": {message}})
	if err != nil {
		log.Fatal(err)
	}
}
