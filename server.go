package main

import (
	"log"
	"myapp/dataloader"
	"myapp/directives"
	"myapp/graph"
	"myapp/graph/generated"
	"myapp/migration"
	"myapp/service"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/go-chi/chi"
)

const defaultPort = "8080"

func main() {
	migration.MigrateTable()

	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	router := chi.NewRouter()
	router.Use(service.CorsMiddleware)
	router.Use(service.Middleware)
	router.Use(dataloader.Middleware)

	c := generated.Config{Resolvers: &graph.Resolver{}}
	c.Directives.IsLogin = directives.IsLogin

	srv := handler.NewDefaultServer(generated.NewExecutableSchema(c))

	router.Handle("/", playground.Handler("GraphQL playground", "/query"))
	router.Handle("/query", srv)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}
