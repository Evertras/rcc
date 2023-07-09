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

func New(cfg Config, coverageRepo CoverageRepository) *Server {
	r := chi.NewRouter()

	// Middleware for all routes
	r.Use(middleware.Logger)
	r.Use(middleware.StripSlashes)

	// Subrouters
	r.Mount("/api", apiRouter(coverageRepo))

	return &Server{
		server: &http.Server{
			Addr:    cfg.Address,
			Handler: r,
		},
	}
}

func apiRouter(coverageRepo CoverageRepository) chi.Router {
	r := chi.NewRouter()

	r.Mount("/v0", v0Router(coverageRepo))

	return r
}

func v0Router(coverageRepo CoverageRepository) chi.Router {
	r := chi.NewRouter()

	// Don't cache values on client side because they will change
	r.Use(middleware.NoCache)

	r.Put("/coverage", v0HandlerCoveragePut(coverageRepo))
	r.Get("/coverage", v0HandlerCoverageGet(coverageRepo))
	r.Get("/badge/coverage", v0HandlerBadgeCoverage(coverageRepo))

	return r
}

// ListenAndServe will run the HTTP server until the given context is
// canceled, or until the underlying HTTP server errors for some other
// reason.
func (s *Server) ListenAndServe(ctx context.Context) error {
	go func() {
		<-ctx.Done()

		log.Println("Context cancelled, shutting down server...")

		timeoutCtx, cancel := context.WithTimeout(context.Background(), time.Second*5)
		defer cancel()

		if err := s.server.Shutdown(timeoutCtx); err != nil {
			log.Println("Error shutting down:", err)
			return
		}

		log.Println("Server shut down successfully")
	}()

	log.Println("Listening at", s.server.Addr)

	return s.server.ListenAndServe()
}

// Handler will return the underlying handler, useful for running a lambda
func (s *Server) Handler() http.Handler {
	return s.server.Handler
}
