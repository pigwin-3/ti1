package main

import (
	"log"
	"ti1/data"
	"ti1/database"
	"ti1/export"
	"time"
)

func main() {
	log.Println("ti1 testing v1.0.1")
	log.Println("Starting...")

	// Setup the database
	err := database.SetupDB()
	if err != nil {
		log.Fatalf("Database setup failed: %v", err)
	}

	// Get the current timestamp
	starttimestamp := time.Now().Format("20060102T150405")
	log.Printf("Starting timestamp: %s", starttimestamp)

	for {
		start := time.Now()

		data, err := data.FetchData(starttimestamp)
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
