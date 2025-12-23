package databaseold

import (
	"context"
	"crypto/md5"
	"database/sql"
	"encoding/hex"
	"fmt"
	"ti1/valki"

	"github.com/valkey-io/valkey-go"
)

func InsertOrUpdateRecordedCall(ctx context.Context, db *sql.DB, values []interface{}, valkeyClient valkey.Client) (int, string, error) {
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

	estimatedVehicleJourneyID := values[0]
	orderID := values[1]
	key := fmt.Sprintf("%v.%v", estimatedVehicleJourneyID, orderID)

	var err error

	// Get the MD5 hash from Valkey
	retrievedHash, err := valki.GetValkeyValue(ctx, valkeyClient, key)
	if err != nil {
		return 0, "", fmt.Errorf("failed to get value from Valkey: %v", err)
	}

	// Check if the retrieved value matches the original MD5 hash
	if retrievedHash != hashString {
		query := `
            INSERT INTO calls (
                estimatedvehiclejourney, "order", stoppointref,
                aimeddeparturetime, expecteddeparturetime,
                aimedarrivaltime, expectedarrivaltime,
                cancellation, actualdeparturetime, actualarrivaltime,
                recorded_data
            )
            VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11)
            ON CONFLICT (estimatedvehiclejourney, "order")
            DO UPDATE SET
                stoppointref = EXCLUDED.stoppointref,
                aimeddeparturetime = EXCLUDED.aimeddeparturetime,
                expecteddeparturetime = EXCLUDED.expecteddeparturetime,
                aimedarrivaltime = EXCLUDED.aimedarrivaltime,
                expectedarrivaltime = EXCLUDED.expectedarrivaltime,
                cancellation = EXCLUDED.cancellation,
                actualdeparturetime = EXCLUDED.actualdeparturetime,
                actualarrivaltime = EXCLUDED.actualarrivaltime,
                recorded_data = EXCLUDED.recorded_data
            RETURNING CASE WHEN xmax = 0 THEN 'insert' ELSE 'update' END, id;
        `
		stmt, err := db.Prepare(query)
		if err != nil {
			return 0, "", fmt.Errorf("error preparing statement: %v", err)
		}
		defer stmt.Close()

		err = valki.SetValkeyValue(ctx, valkeyClient, key, hashString)
		if err != nil {
			return 0, "", fmt.Errorf("failed to set value in Valkey: %v", err)
		}

		var action string
		var id int
		err = stmt.QueryRow(values...).Scan(&action, &id)
		if err != nil {
			return 0, "", fmt.Errorf("error executing statement: %v", err)
		}
		return id, action, nil
	} else {
		return 0, "none", nil
	}
}
