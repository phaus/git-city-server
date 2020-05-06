package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/grandcat/zeroconf"
	"github.com/julienschmidt/httprouter"
)

// VersionString - the version of the application
var VersionString string

const (
	port = 3000
)

func main() {
	// Setup our service export
	host, _ := os.Hostname()

	server, err := zeroconf.Register("git-city service", "_git-city._tcp", "local.", port, []string{"txtv=1", fmt.Sprintf("host=%s", host)}, nil)
	if err != nil {
		panic(err)
	}
	defer server.Shutdown()

	fmt.Printf("\nrunning on http://%s:%d\n", host, port)

	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	// Clean exit.
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)
	select {
	case <-sig:
		// Exit by user
	case <-time.After(time.Second * 120):
		// Exit by timeout
	}
}

func run() error {
	router := httprouter.New()
	router.GET("/", foo)
	return http.ListenAndServe(fmt.Sprintf(":%d", port), router)
}

func foo(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	payload, err := createPayload()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	js, err := json.Marshal(payload)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Server", "git-city-server")
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func createPayload() ([]map[string]interface{}, error) {
	payloads := make([]map[string]interface{}, 0)
	for i := 1; i < 82; i++ {
		payload := make(map[string]interface{})
		j := rand.Intn(100) * 100
		payload[fmt.Sprintf("entry-%d", i)] = j
		payloads = append(payloads, payload)
	}
	return payloads, nil
}
