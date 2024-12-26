package main

import (
	"log"
	"ti1/config"
	"ti1/data"
	"ti1/export"
)

func main() {
	config.PrintDBConfig()

	data, err := data.FetchData()
	if err != nil {
		log.Fatal(err)
	}

	//export.ExportToCSV(data)
	export.DBData(data)
	//export.PrintData(data)

	//log.Printf("Data fetched successfully: %+v", data)

	//export.PrintData(data)
}
