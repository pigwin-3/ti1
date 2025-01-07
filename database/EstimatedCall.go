package database

import (
	"crypto/md5"
	"database/sql"
	"encoding/hex"
	"fmt"
	"log"
	"ti1/config"
	"ti1/valkey"
)

func InsertOrUpdateEstimatedCall(db *sql.DB, values []interface{}) (int, string, error) {
	// Replace empty strings with nil for timestamp fields
	for i, v := range values {
		if str, ok := v.(string); ok && str == "" {
			values[i] = nil
		}
	}

	// Convert values to a single string and hash it using MD5
	var valuesString string
	for _, v := range values {
		if v != nil {
			valuesString += fmt.Sprintf("%v", v)
		}
	}
	hash := md5.Sum([]byte(valuesString))
	hashString := hex.EncodeToString(hash[:])
	println(hashString)

	estimatedVehicleJourneyID := values[0]
	orderID := values[1]
	key := fmt.Sprintf("%v.%v", estimatedVehicleJourneyID, orderID)
	fmt.Printf("Estimated Vehicle Journey ID: %v, Order ID: %v\n", estimatedVehicleJourneyID, orderID)

	// Connect to Valkey
	valkeyClient, err := config.ConnectToValkey("config/conf.json")
	if err != nil {
		log.Fatalf("Failed to connect to Valkey: %v", err)
	}
	defer config.DisconnectFromValkey(valkeyClient)

	// Set the MD5 hash in Valkey
	err = valkey.SetValkeyValue(valkeyClient, key, hashString)
	if err != nil {
		log.Fatalf("Failed to set value in Valkey: %v", err)
	}

	// Get the MD5 hash from Valkey
	retrievedHash, err := valkey.GetValkeyValue(valkeyClient, key)
	if err != nil {
		log.Fatalf("Failed to get value from Valkey: %v", err)
	}

	// Check if the retrieved value matches the original MD5 hash
	if retrievedHash != hashString {
		log.Fatalf("Retrieved hash does not match the original hash. Original: %s, Retrieved: %s", hashString, retrievedHash)
	} else {
		fmt.Println("Retrieved hash matches the original hash.")
	}

	query := `
        INSERT INTO calls (
            estimatedvehiclejourney, "order", stoppointref,
            aimeddeparturetime, expecteddeparturetime,
            aimedarrivaltime, expectedarrivaltime,
            cancellation, estimated_data
        )
        VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)
        ON CONFLICT (estimatedvehiclejourney, "order")
        DO UPDATE SET
            stoppointref = EXCLUDED.stoppointref,
            aimeddeparturetime = EXCLUDED.aimeddeparturetime,
            expecteddeparturetime = EXCLUDED.expecteddeparturetime,
            aimedarrivaltime = EXCLUDED.aimedarrivaltime,
            expectedarrivaltime = EXCLUDED.expectedarrivaltime,
            cancellation = EXCLUDED.cancellation,
            estimated_data = EXCLUDED.estimated_data
        RETURNING CASE WHEN xmax = 0 THEN 'insert' ELSE 'update' END, id;
    `
	stmt, err := db.Prepare(query)
	if err != nil {
		return 0, "", fmt.Errorf("error preparing statement: %v", err)
	}
	defer stmt.Close()

	var action string
	var id int
	err = stmt.QueryRow(values...).Scan(&action, &id)
	if err != nil {
		if 1 == 0 {
			fmt.Println("Executing query:", query)
			for i, v := range values {
				fmt.Printf("Value %d: (%v)\n", i+1, v)
			}

		}
		return 0, "", fmt.Errorf("error executing statement: %v", err)
	}
	return id, action, nil
}
