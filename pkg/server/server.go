package server

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type Server struct {
	server *http.Server
}

func New() *Server {
	r := chi.NewRouter()

	// Middleware for all routes
	r.Use(middleware.Logger)

	// Subrouters
	r.Mount("/api", apiRouter())

	return &Server{
		server: &http.Server{
			Handler: r,
		},
	}
}

func apiRouter() chi.Router {
	r := chi.NewRouter()

	r.Mount("/v0", v0Router())

	return r
}

func v0Router() chi.Router {
	r := chi.NewRouter()

	// Don't cache values on client side because they may change
	r.Use(middleware.NoCache)

	r.Put("/coverage/value100", v0HandlerCoveragePut())

	return r
}

// ListenAndServe will run the HTTP server until the given context is
// canceled, or until the underlying HTTP server errors for some other
// reason.
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

// Handle will handle the given response/request directly, useful for lambdas
func (s *Server) Handle(w http.ResponseWriter, r *http.Request) {
	s.server.Handler.ServeHTTP(w, r)
}
