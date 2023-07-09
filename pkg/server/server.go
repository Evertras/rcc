package server

import (
	"context"
	"log"
	"net/http"
	"time"
)

type Server struct {
	server *http.Server
}

func New() *Server {
	return &Server{
		server: &http.Server{},
	}
}

func (s *Server) ListenAndServe(ctx context.Context) error {
	go func() {
		<-ctx.Done()

		timeoutCtx, cancel := context.WithTimeout(context.Background(), time.Second*5)
		defer cancel()

		if err := s.server.Shutdown(timeoutCtx); err != nil {
			log.Println("Error shutting down:", err)
		}
	}()

	return s.server.ListenAndServe()
}
