package http

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/solardome/gamepulse-platform/accounts-gateway/graph"
	"github.com/solardome/gamepulse-platform/accounts-gateway/graph/generated"
)

func NewRouter(logger *slog.Logger, resolver *graph.Resolver) http.Handler {
	graphServer := handler.New(
		generated.NewExecutableSchema(generated.Config{
			Resolvers: resolver,
		}),
	)

	graphServer.AddTransport(transport.Options{})
	graphServer.AddTransport(transport.GET{})
	graphServer.AddTransport(transport.POST{})
	graphServer.AddTransport(transport.MultipartForm{})

	graphServer.Use(extension.Introspection{})

	mux := http.NewServeMux()
	mux.Handle("/", playground.Handler("GamePulse GraphQL", "/query"))
	mux.Handle("/query", graphServer)
	mux.HandleFunc("/livez", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})

	return logRequests(logger, mux)
}

func logRequests(logger *slog.Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		logger.Info(
			"http request",
			"method", r.Method,
			"path", r.URL.Path,
			"duration_ms", time.Since(start).Milliseconds(),
		)
	})
}
