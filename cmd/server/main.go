package main

import (
	"log"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/lru"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/eldhosereji541/task-tracker-backend/internal/auth"
	"github.com/eldhosereji541/task-tracker-backend/internal/graph"
	"github.com/eldhosereji541/task-tracker-backend/internal/store"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/vektah/gqlparser/v2/ast"
)

const defaultPort = "8080"

var cors = []string{"*"} // todo: restrict this in production to specific origins

func main() {
	// ========== Echo Server and Middleware Setup========== //
	e := echo.New()
	e.Use(middleware.BodyLimit("32M"))
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     cors,
		AllowMethods:     []string{echo.GET, echo.POST, echo.PUT, echo.DELETE, echo.OPTIONS},
		AllowCredentials: true,
	}))
	e.Use(middleware.Secure())

	// graphql setup
	repo := store.NewStore()
	srv := handler.New(graph.NewExecutableSchema(graph.Config{
		Resolvers: &graph.Resolver{
			TaskRepo: repo,
		},
	}))

	srv.AddTransport(transport.Options{})
	srv.AddTransport(transport.GET{})
	srv.AddTransport(transport.POST{})

	srv.SetQueryCache(lru.New[*ast.QueryDocument](1000))

	srv.Use(extension.Introspection{})
	srv.Use(extension.AutomaticPersistedQuery{
		Cache: lru.New[string](100),
	})

	// GraphQL playground handler for testing and development.
	e.GET("/", echo.WrapHandler(playground.Handler("GraphQL playground", "/query")))

	// GraphQL endpoint with auth middleware
	e.POST("/query", echo.WrapHandler(auth.AuthMiddleware(srv)))
	e.GET("/query", echo.WrapHandler(auth.AuthMiddleware(srv)))

	err := godotenv.Load()
	if err != nil {
		// fallback path
		err = godotenv.Load("../../.env")
		if err != nil {
			log.Println("No .env file found")
		}
	}

	// ========== Start Server ========== //
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	// jwt secret fetch and set
	secret := os.Getenv("JWT_SECRET")
	if len(secret) == 0 || len(secret) < 32 {
		log.Fatal("JWT_SECRET environment variable is not set or invalid")
	}
	auth.JwtSecret = []byte(secret)

	// todo: add graceful shutdown
	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	e.Logger.Fatal(e.Start(":" + port))
}
