package main

import (
	"log"
	"ti1/data"
	"ti1/database"
	"ti1/export"
	"time"
)

func main() {
	log.Println("Starting...")

	// Setup the database
	err := database.SetupDB()
	if err != nil {
		log.Fatalf("Database setup failed: %v", err)
	}

	for {
		start := time.Now()

		data, err := data.FetchData()
		if err != nil {
			log.Fatal(err)
		}

		export.DBData(data)

		log.Println("finished in", time.Since(start))
		elapsed := time.Since(start)
		if elapsed < 5*time.Minute {
			log.Printf("starting again in %v", 5*time.Minute-elapsed)
			time.Sleep(1*time.Minute - elapsed)
		}
	}
}
