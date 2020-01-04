package hello

import (
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/go-redis/redis"
	"github.com/sirupsen/logrus"
)

type greeter struct {
	version string
	greetee string
}

var redisClient *redis.Client

var log = logrus.New()

func init() {
	log.Formatter = &logrus.JSONFormatter{}
	log.Level = logrus.DebugLevel
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

func (g *greeter) loggingPlaygroundHandler(w http.ResponseWriter, r *http.Request) {
	httplog := httpLogger(log, r)

	level := r.URL.Query().Get("level")
	message := r.URL.Query().Get("message")

	if level == "info" {
		httplog.Info(message)
	} else if level == "warn" {
		httplog.Warn(message)
	} else if level == "error" {
		httplog.Error(message)
	}

	fmt.Fprintf(w, "Logged: Level: %s, Message: %s\n", level, message)
}

func httpLogger(logger *logrus.Logger, r *http.Request) *logrus.Entry {
	return logger.WithFields(logrus.Fields{
		"httpRequest": map[string]interface{}{
			"method":    r.Method,
			"url":       r.URL.String(),
			"userAgent": r.Header.Get("User-Agent"),
			"referrer":  r.Header.Get("Referer"),
		},
	})
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
	log.Info("Serving health endpoints on :8081")
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

	helloRouter := http.NewServeMux()
	helloRouter.HandleFunc("/", g.rootHandler)
	helloRouter.HandleFunc("/log", g.loggingPlaygroundHandler)
	helloServer := http.Server{
		Addr:    ":8080",
		Handler: helloRouter,
	}
	log.Info("Serving on :8080")
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
