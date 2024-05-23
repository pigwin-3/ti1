package main

import (
	"log"
	"ti1/data" // Import the data package
	"ti1/export"
)

func main() {
	data, err := data.FetchData() // Use the FetchData function from the data package
	if err != nil {
		log.Fatal(err)
	}

	//log.Printf("Data fetched successfully: %+v", data)

	export.PrintData(data)
}
