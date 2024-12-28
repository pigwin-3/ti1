package main

import (
	"log"
	"ti1/data"
	"ti1/export"
)

func main() {
	//config.PrintDBConfig()

	for i := 0; i < 10; i++ {
		data, err := data.FetchData()
		if err != nil {
			log.Fatal(err)
		}

		//export.ExportToCSV(data)
		export.DBData(data)
	}
	println(":)")
	//export.PrintData(data)

	//log.Printf("Data fetched successfully: %+v", data)

	//export.PrintData(data)

}
