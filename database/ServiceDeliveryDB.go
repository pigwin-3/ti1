package database

import (
	"database/sql"
	"fmt"
)

func InsertServiceDelivery(db *sql.DB) (int, error) {
	fmt.Println("Inserting ServiceDelivery...")
	var id int

	responseTimestamp := "2024-05-25T10:24:29.353864654+02:00"
	recordedAtTime := "2024-05-25T10:24:29.353864654+02:00"

	err := db.QueryRow("INSERT INTO public.ServiceDelivery (ResponseTimestamp, RecordedAtTime) VALUES ($1, $2) RETURNING ID", responseTimestamp, recordedAtTime).Scan(&id)
	fmt.Println(err)
	if err != nil {
		return 0, err
	}
	fmt.Println("ID:", id)
	fmt.Println("ServiceDelivery inserted successfully!")
	return id, nil
}
