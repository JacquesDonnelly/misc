package hello

import (
	"fmt"
	"net/http"
	"time"
)

var (
	healthy  = true
	ready    = false
	stopping = false
)

func setUnhealthyHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Going unhealthy")
	healthy = false
}
func setHealthyHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Going healthy")
	healthy = true
}

func setUnreadyHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Going unready")
	ready = false
}
func setReadyHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Going ready")
	ready = true
}

func healthzHandler(w http.ResponseWriter, r *http.Request) {
	if healthy {
		w.WriteHeader(200)
		w.Write([]byte("ok"))
	} else {
		w.WriteHeader(500)
		w.Write([]byte("error"))
	}
}

func readyHandler(w http.ResponseWriter, r *http.Request) {
	if ready {
		w.WriteHeader(200)
		w.Write([]byte("ok"))
	} else {
		w.WriteHeader(500)
		w.Write([]byte(fmt.Sprintf("error")))
	}
}

func delayReady() {
	t := time.NewTimer(10 * time.Second)
	<-t.C
	ready = true
}

func preStopHandler(w http.ResponseWriter, r *http.Request) {
	stopping = true

	// perform some shutdown work - k8s waits
	t := time.NewTimer(10 * time.Second)
	<-t.C

	// return ok
	w.WriteHeader(200)
	w.Write([]byte("ready to stop, kubernetes, kill me now"))
}

func stoppingStateHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
	if stopping {
		fmt.Fprintf(w, "Stopping!")
	} else {
		fmt.Fprintf(w, "Not stopping yet")
	}
}
