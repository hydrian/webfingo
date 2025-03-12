package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"

	"webfingo/internal/webfingo"
)

func main() {
	logger := webfingo.DefaultLogger
	conf := getConf(logger)

	// Create database instance with connection handling moved to NewDatabase
	dbInstance, err := webfingo.NewDatabase(conf.GetDBConnectionString())
	if err != nil {
		logger.Fatalf("Unable to connect to database: %v\n", err)
	}
	defer dbInstance.Close()

	logger.Println("Successfully connected to database")

	// Define WebFinger handler
	http.HandleFunc("/.well-known/webfinger", func(w http.ResponseWriter, r *http.Request) {
		webfingo.HandleWebfingerRequest(w, r, dbInstance, conf.Keycloak, logger)
	})

	// Get webserver configuration
	webserverConfig := conf.GetWebfingoWebserverConfig()
	port := webserverConfig.Port
	serverAddr := fmt.Sprintf(":%d", port)

	// Start HTTP server
	logger.Printf("Server starting on port %d...", port)
	if err := http.ListenAndServe(serverAddr, nil); err != nil {
		logger.Fatalf("Server failed to start: %v", err)
	}
}

func getConf(logger *webfingo.Logger) *webfingo.Config {
	// Parse command line flags
	configPath := flag.String("config", "", "Path to configuration file (required)")
	flag.Parse()

	// Check if config flag is provided
	if *configPath == "" {
		flag.Usage()
		logger.Println("Error: config flag is required")
		os.Exit(1)
	}

	// Load configuration
	c, err := webfingo.LoadConfig(*configPath)
	if err != nil {
		logger.Fatalf("Failed to load configuration: %v", err)
	}

	return c
}
