package database

import (
	"database/sql"
	"fmt"
	"ti1/config"
)

func GetDatasetVariable(config config.Config) string {
	if config.DatasetId != "" {
		fmt.Println(config.DatasetId)
		return config.DatasetId
	} else if config.ExcludedDatasetIds != "" {
		result := "EX." + config.ExcludedDatasetIds
		fmt.Println(result)
		return result
	}
	fmt.Println("")
	return ""
}

func InsertServiceDelivery(db *sql.DB, responseTimestamp string, recordedAtTime string) (int, error) {
	fmt.Println("Inserting ServiceDelivery...")
	var id int

	// Load configuration
	config, err := config.LoadConfig()
	if err != nil {
		fmt.Println("Error loading config:", err)
		return 0, err
	}

	// Get dataset variable
	datasetVariable := GetDatasetVariable(config)

	err = db.QueryRow("INSERT INTO public.ServiceDelivery (ResponseTimestamp, RecordedAtTime, DatasetVariable) VALUES ($1, $2, $3) RETURNING ID", responseTimestamp, recordedAtTime, datasetVariable).Scan(&id)
	if err != nil {
		fmt.Println(err)
		return 0, err
	}
	//fmt.Println("ServiceDelivery inserted successfully! (", id, ")")
	return id, nil
}

func UpdateServiceDeliveryData(db *sql.DB, id int, data string) error {
	fmt.Println("Updating ServiceDelivery data...")
	_, err := db.Exec("UPDATE public.ServiceDelivery SET Data = $1 WHERE ID = $2", data, id)
	if err != nil {
		fmt.Println(err)
		return err
	}
	fmt.Println("Finished with this ServiceDelivery!")
	return nil
}
