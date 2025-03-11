package main

import (
	"context"
	"database/sql"
	"flag"
	"log"
	"net/http"
	"os"

	"webfingo/internal/webfingo"

	_ "github.com/lib/pq"
)

func main() {

	conf := getConf()

	// Connect to PostgreSQL using the configuration
	db, err := sql.Open("postgres", conf.GetDBConnectionString())
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}
	defer db.Close()

	// Verify database connection
	if err := db.PingContext(context.Background()); err != nil {
		log.Fatalf("Unable to ping database: %v\n", err)
	}
	log.Println("Successfully connected to database")

	// Create database instance
	dbInstance := webfingo.NewDatabase(db)

	// Define WebFinger handler
	http.HandleFunc("/.well-known/webfinger", func(w http.ResponseWriter, r *http.Request) {
		webfingo.HandleWebfingerRequest(w, r, dbInstance, conf)
	})

	// Start HTTP server
	log.Println("Server starting on port 8080...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}

func getConf() *webfingo.Config {
	// Parse command line flags
	configPath := flag.String("config", "", "Path to configuration file (required)")
	flag.Parse()

	// Check if config flag is provided
	if *configPath == "" {
		flag.Usage()
		log.Println("Error: config flag is required")
		os.Exit(1)
	}

	// Load configuration
	c, err := webfingo.LoadConfig(*configPath)
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	return c
}
