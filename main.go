package main

import (
	"log"
	"ti1/config"
	"ti1/data"
	"ti1/database"
	"ti1/export"
	"time"
)

func main() {
	log.Println("ti1 v0.2.1")
	log.Println("Starting...")

	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Setup the database
	err = database.SetupDB()
	if err != nil {
		log.Fatalf("Database setup failed: %v", err)
	}

	// Get the current timestamp
	starttimestamp := time.Now().Format("20060102T150405")
	log.Printf("Starting timestamp: %s", starttimestamp)

	for {
		start := time.Now()

		data, err := data.FetchData(starttimestamp, cfg.DatasetId, cfg.ExcludedDatasetIds)
		if err != nil {
			log.Fatal(err)
		}

		export.DBData(data)

		log.Println("finished in", time.Since(start))
		elapsed := time.Since(start)
		if elapsed < 5*time.Minute {
			log.Printf("starting again in %v", 5*time.Minute-elapsed)
			time.Sleep(5*time.Minute - elapsed)
		}
	}
}
