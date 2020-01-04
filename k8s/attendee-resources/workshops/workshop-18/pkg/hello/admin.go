package hello

import (
	"fmt"
	"net/http"
	"time"
)

var (
	healthy = true
	ready   = false
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
	t := time.NewTimer(1 * time.Second)
	<-t.C
	ready = true
}
