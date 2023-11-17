package main

import (
	"encoding/json"
	"flag"
	"log"
	"net/http"
	"os"
	"runtime"

	"gopkg.in/yaml.v2"
)

func main() {

	configFile := flag.String("config", "./config.yaml", "Path to the configuration file")
	flag.Parse()

	// Reading the configuration file
	config = Config{}
	data, err := os.ReadFile(*configFile)
	if err != nil {
		log.Fatalf("Error reading config file: %v", err)
	}

	err = yaml.Unmarshal(data, &config)
	if err != nil {
		log.Fatalf("Error parsing config file: %v", err)
	}

	// Set the OS variable
	config.OS = runtime.GOOS

	http.HandleFunc("/search", searchHandler)
	log.Println("Starting server on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func searchHandler(w http.ResponseWriter, r *http.Request) {
	// Extract query parameter
	query := r.URL.Query().Get("q")
	if query == "" {
		http.Error(w, "Query parameter 'q' is required", http.StatusBadRequest)
		return
	}

	// Perform the search
	results, err := performSearch(query)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Respond with JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(results)
}
