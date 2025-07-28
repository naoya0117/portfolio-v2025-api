package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gorilla/mux"
	"github.com/rs/cors"

	"github.com/naoya0117/portfolio-v2025-api/internal/auth"
	"github.com/naoya0117/portfolio-v2025-api/internal/database"
	"github.com/naoya0117/portfolio-v2025-api/internal/generated"
	"github.com/naoya0117/portfolio-v2025-api/internal/resolvers"
)

const defaultPort = "8080"

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	// Initialize database connection (optional)
	db, err := database.NewConnection()
	if err != nil {
		log.Printf("Warning: Failed to connect to database: %v. Using mock data.", err)
		db = nil
	}
	
	if db != nil {
		log.Println("Successfully connected to database")
		
		// Create tables if they don't exist
		if err := db.CreateTables(); err != nil {
			log.Fatalf("Failed to create tables: %v", err)
		}
		
		// Run migrations to add like_count columns
		if err := db.MigrateTables(); err != nil {
			log.Fatalf("Failed to migrate tables: %v", err)
		}
		
		log.Println("Database tables initialized and migrated")
		
		// Seed initial data (commented out to start with empty database)
		// if err := db.SeedData(); err != nil {
		//	log.Printf("Warning: Failed to seed data: %v", err)
		// }
		
		defer db.Close()
	} else {
		log.Println("Using mock data (no database connection)")
	}

	// Initialize resolver with database connection
	resolver := &resolvers.Resolver{DB: db}

	srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: resolver}))

	// Create router with debug logging
	router := mux.NewRouter()
	
	// Auth endpoints (no auth middleware)
	router.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("[ROUTER] Login handler called: %s %s\n", r.Method, r.URL.Path)
		auth.LoginHandler(w, r)
	}).Methods("POST", "OPTIONS")
	
	// Public endpoints (no auth required)
	// Only enable playground in development
	if os.Getenv("GO_ENV") != "production" {
		router.Handle("/", playground.Handler("GraphQL playground", "/query"))
	}
	router.Handle("/query", srv)
	
	// Protected admin endpoints
	// Only enable admin playground in development
	if os.Getenv("GO_ENV") != "production" {
		router.Handle("/admin", auth.AuthMiddleware(playground.Handler("GraphQL playground (Admin)", "/admin/query")))
	}
	router.Handle("/admin/query", auth.AuthMiddleware(srv))

	// Get allowed origins from environment or use defaults
	allowedOrigins := []string{
		"http://localhost:3000", 
		"http://localhost:8000",
		"http://localhost:3001",
	}
	
	// Add custom origins from environment if specified
	if customOrigins := os.Getenv("CORS_ALLOWED_ORIGINS"); customOrigins != "" {
		origins := strings.Split(customOrigins, ",")
		for _, origin := range origins {
			allowedOrigins = append(allowedOrigins, strings.TrimSpace(origin))
		}
	}

	c := cors.New(cors.Options{
		AllowedOrigins:   allowedOrigins,
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
		Debug:            os.Getenv("GO_ENV") == "development",
	})

	// Add global request logging
	finalHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("[SERVER] %s %s from %s\n", r.Method, r.URL.Path, r.RemoteAddr)
		c.Handler(router).ServeHTTP(w, r)
	})

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, finalHandler))
}