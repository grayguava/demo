package dash

import (
	"fmt"
	"net/http"
	"time"
)

var port = 3847

type Server struct {
	addr   string
	svc    interface {
		IsRunning() bool
		Start()
		Stop()
		LastSyncTime() time.Time
		GetConfig() interface{}
	}
}

func New(svc interface {
	IsRunning() bool
	Start()
	Stop()
	LastSyncTime() time.Time
	GetConfig() interface{}
}) *Server {
	return &Server{
		addr: fmt.Sprintf(":%d", port),
		svc:  svc,
	}
}

func (s *Server) Port() int {
	return port
}

func (s *Server) Start() error {
	mux := http.NewServeMux()
	mux.HandleFunc("/", handleIndex)
	mux.HandleFunc("/api/status", handleStatus(s.svc))
	mux.HandleFunc("/api/start", handleStart(s.svc))
	mux.HandleFunc("/api/stop", handleStop(s.svc))
	mux.HandleFunc("/api/save", handleSave(s.svc))

	server := &http.Server{
		Addr:         s.addr,
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	fmt.Printf("Dashboard starting on http://localhost:%d\n", port)
	return server.ListenAndServe()
}