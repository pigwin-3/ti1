package database

import (
	"database/sql"
	"fmt"
)

func InsertServiceDelivery(db *sql.DB, responseTimestamp string, recordedAtTime string) (int, error) {
	fmt.Println("Inserting ServiceDelivery...")
	var id int

	err := db.QueryRow("INSERT INTO public.ServiceDelivery (ResponseTimestamp, RecordedAtTime) VALUES ($1, $2) RETURNING ID", responseTimestamp, recordedAtTime).Scan(&id)
	fmt.Println(err)
	if err != nil {
		return 0, err
	}
	fmt.Println("ID:", id)
	fmt.Println("ServiceDelivery inserted successfully!")
	return id, nil
}
