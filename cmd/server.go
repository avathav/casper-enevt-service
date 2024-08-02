package cmd

import (
	"net/http"
	_ "net/http/pprof"

	"event-service/graph"
	"event-service/internal/config"
	"event-service/internal/di"
	ihtttp "event-service/internal/http"
	"event-service/internal/profiler"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

const defaultPort = "8080"

var apiCmd = &cobra.Command{
	Use:   "start-api",
	Short: "cli that starts events service",
	Run:   startApi,
}

func init() {
	rootCmd.AddCommand(apiCmd)

}

func startApi(cmd *cobra.Command, _ []string) {
	defer di.CloseAllExchangeConnections()

	r := chi.NewRouter()
	r.Use(
		middleware.RequestID,
		middleware.RealIP,
		middleware.Logger,
		middleware.Recoverer,
		ihtttp.GORMConnectionMiddleware(di.GORM()),
	)

	r.Mount("/debug/pprof", profiler.Router())

	port := config.GetStringOrFallback("API.PORT", defaultPort)

	resolver, err := di.DefaultGraphQLApiResolver()
	if err != nil {
		log.WithError(err).Panic("cannot create graphql resolver")
	}

	srv := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: resolver}))

	r.Handle("/", playground.Handler("GraphQL playground", "/query"))
	r.Handle("/query", srv)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
