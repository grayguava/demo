package main

import (
	"fmt"
	"net"
	"net/http"
	"time"
)

const defaultPort = 3847

func startServer(state *AppState) (net.Listener, error) {
	mux := http.NewServeMux()
	mux.HandleFunc("/", handleIndex)
	mux.HandleFunc("/api/config", handleConfig(state))
	mux.HandleFunc("/api/sync", handleSync(state))
	mux.HandleFunc("/api/status", handleStatus(state))

	ln, err := net.Listen("tcp", fmt.Sprintf(":%d", defaultPort))
	if err != nil {
		return nil, err
	}

	srv := &http.Server{
		Handler:      mux,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 60 * time.Second,
	}

	go srv.Serve(ln)
	return ln, nil
}