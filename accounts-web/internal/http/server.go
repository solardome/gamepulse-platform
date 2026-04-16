package http

import (
	"context"
	"embed"
	"html/template"
	"log/slog"
	"net/http"
	"time"

	"github.com/solardome/gamepulse-platform/accounts-web/internal/graphql"
)

//go:embed templates/*.tmpl
var templateFS embed.FS

type Server struct {
	logger    *slog.Logger
	graphql   *graphql.Client
	templates *template.Template
}

type resultViewData struct {
	Message string
	Service string
	Error   string
}

func NewServer(logger *slog.Logger, graphQLClient *graphql.Client) (*Server, error) {
	templates, err := template.ParseFS(templateFS, "templates/*.tmpl")
	if err != nil {
		return nil, err
	}

	return &Server{
		logger:    logger,
		graphql:   graphQLClient,
		templates: templates,
	}, nil
}

func (s *Server) Routes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /", s.handleIndex)
	mux.HandleFunc("POST /actions/ping", s.handlePing)
	mux.HandleFunc("GET /livez", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})

	return s.logRequests(mux)
}

func (s *Server) handleIndex(w http.ResponseWriter, _ *http.Request) {
	data := struct {
		Title string
	}{
		Title: "GamePulse Accounts",
	}

	if err := s.templates.ExecuteTemplate(w, "index", data); err != nil {
		http.Error(w, "template error", http.StatusInternalServerError)
	}
}

func (s *Server) handlePing(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 4*time.Second)
	defer cancel()

	result, err := s.graphql.Ping(ctx)
	if err != nil {
		_ = s.templates.ExecuteTemplate(w, "result", resultViewData{
			Error: err.Error(),
		})
		return
	}

	if err := s.templates.ExecuteTemplate(w, "result", resultViewData{
		Message: result.Message,
		Service: result.Service,
	}); err != nil {
		http.Error(w, "template error", http.StatusInternalServerError)
	}
}

func (s *Server) logRequests(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		s.logger.Info(
			"http request",
			"method", r.Method,
			"path", r.URL.Path,
			"duration_ms", time.Since(start).Milliseconds(),
		)
	})
}
