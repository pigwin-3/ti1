package database

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
	fmt.Println("ServiceDelivery inserted successfully! (", id, ")")
	return id, nil
}
