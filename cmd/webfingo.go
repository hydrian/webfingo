package main

import (
	"flag"
	"log"
	"net/http"
	"os"

	"webfingo/internal/webfingo"

	_ "github.com/lib/pq"
)

func main() {
	conf := getConf()

	// Create database instance with connection handling moved to NewDatabase
	dbInstance, err := webfingo.NewDatabase(conf.GetDBConnectionString())
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}
	defer dbInstance.Close()

	log.Println("Successfully connected to database")

	// Define WebFinger handler
	http.HandleFunc("/.well-known/webfinger", func(w http.ResponseWriter, r *http.Request) {
		webfingo.HandleWebfingerRequest(w, r, dbInstance, conf.Keycloak)
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
