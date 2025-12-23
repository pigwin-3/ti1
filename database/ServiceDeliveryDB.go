package database

import (
	"database/sql"
	"fmt"
)

func InsertServiceDelivery(db *sql.DB, responseTimestamp string, recordedAtTime string) (int, error) {
	var id int
	err := db.QueryRow("INSERT INTO public.ServiceDelivery (ResponseTimestamp, RecordedAtTime) VALUES ($1, $2) RETURNING ID", responseTimestamp, recordedAtTime).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("failed to insert service delivery: %w", err)
	}
	return id, nil
}

func UpdateServiceDeliveryData(db *sql.DB, id int, data string) error {
	_, err := db.Exec("UPDATE public.ServiceDelivery SET Data = $1 WHERE ID = $2", data, id)
	if err != nil {
		return fmt.Errorf("failed to update service delivery data: %w", err)
	}
	return nil
}
