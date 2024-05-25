package main

import (
	"log"
	"ti1/config"
	"ti1/data"
	"ti1/database"
	"ti1/export"
)

func main() {
	config.PrintDBConfig()
	config.ConnectToPostgreSQL()

	db, err := config.ConnectToPostgreSQL()
	if err != nil {
		log.Fatal(err)
	}
	//log.Printf("DB: %+v", db)

	database.InsertServiceDelivery(db)

	data, err := data.FetchData()
	if err != nil {
		log.Fatal(err)
	}

	//log.Printf("Data fetched successfully: %+v", data)

	export.PrintData(data)
}
