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
	// Load environment variables before anything else
	if err := godotenv.Load(); err != nil {
		if err = godotenv.Load("../../.env"); err != nil {
			log.Println("No .env file found, relying on environment variables")
		}
	}

	secret := os.Getenv("JWT_SECRET")
	if len(secret) < 32 {
		log.Fatal("JWT_SECRET must be at least 32 characters")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	tokenSvc := auth.NewTokenService([]byte(secret))

	// ========== Echo Server and Middleware Setup ========== //
	e := echo.New()
	e.Use(middleware.BodyLimit("32M"))
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     cors,
		AllowMethods:     []string{echo.GET, echo.POST, echo.PUT, echo.DELETE, echo.OPTIONS},
		AllowCredentials: true,
	}))
	e.Use(middleware.Secure())

	repo := store.NewStore()
	srv := handler.New(graph.NewExecutableSchema(graph.Config{
		Resolvers: &graph.Resolver{
			TaskRepo: repo,
			TokenSvc: tokenSvc,
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

	authMw := auth.AuthMiddleware(tokenSvc)

	e.GET("/", echo.WrapHandler(playground.Handler("GraphQL playground", "/query")))
	e.POST("/query", echo.WrapHandler(authMw(srv)))
	e.GET("/query", echo.WrapHandler(authMw(srv)))

	// todo: add graceful shutdown
	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	e.Logger.Fatal(e.Start(":" + port))
}
