package databaseold

import (
	"database/sql"
	"fmt"
)

func InsertServiceDelivery(db *sql.DB, responseTimestamp string, recordedAtTime string) (int, error) {
	fmt.Println("Inserting ServiceDelivery...")
	var id int

	err := db.QueryRow("INSERT INTO public.ServiceDelivery (ResponseTimestamp, RecordedAtTime) VALUES ($1, $2) RETURNING ID", responseTimestamp, recordedAtTime).Scan(&id)
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
