package database

import (
	"context"
	"crypto/md5"
	"database/sql"
	"encoding/hex"
	"fmt"
	"ti1/valki"

	"github.com/valkey-io/valkey-go"
)

func InsertOrUpdateEstimatedVehicleJourney(ctx context.Context, db *sql.DB, values []interface{}, valkeyClient valkey.Client) (int, string, error) {
	// Generate a key using lineref, directionref, datasource, and datedvehiclejourneyref
	lineref := values[2]
	directionref := values[3]
	datasource := values[4]
	datedvehiclejourneyref := values[5]
	key := fmt.Sprintf("%v.%v.%v.%v", lineref, directionref, datasource, datedvehiclejourneyref)

	// Convert values to a single string and hash it using MD5
	var valuesString string
	for _, v := range values {
		if v != nil {
			valuesString += fmt.Sprintf("%v", v)
		}
	}
	hash := md5.Sum([]byte(valuesString))
	hashString := hex.EncodeToString(hash[:])

	// Get the MD5 hash from Valkey
	retrievedHash, err := valki.GetValkeyValue(ctx, valkeyClient, key)
	if err != nil {
		return 0, "", fmt.Errorf("failed to get value from Valkey: %v", err)
	}

	// Check if the retrieved value matches the original MD5 hash
	if retrievedHash != hashString {
		query := `
		INSERT INTO estimatedvehiclejourney (servicedelivery, recordedattime, lineref, directionref, datasource, datedvehiclejourneyref, vehiclemode, dataframeref, originref, destinationref, operatorref, vehicleref, cancellation, other, firstservicedelivery)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $1)
		ON CONFLICT (lineref, directionref, datasource, datedvehiclejourneyref)
		DO UPDATE SET
			servicedelivery = EXCLUDED.servicedelivery,
			recordedattime = EXCLUDED.recordedattime,
			vehiclemode = COALESCE(EXCLUDED.vehiclemode, estimatedvehiclejourney.vehiclemode),
			dataframeref = COALESCE(EXCLUDED.dataframeref, estimatedvehiclejourney.dataframeref),
			originref = COALESCE(EXCLUDED.originref, estimatedvehiclejourney.originref),
			destinationref = COALESCE(EXCLUDED.destinationref, estimatedvehiclejourney.destinationref),
			operatorref = COALESCE(EXCLUDED.operatorref, estimatedvehiclejourney.operatorref),
			vehicleref = COALESCE(EXCLUDED.vehicleref, estimatedvehiclejourney.vehicleref),
			cancellation = COALESCE(EXCLUDED.cancellation, estimatedvehiclejourney.cancellation),
			other = COALESCE(EXCLUDED.other, estimatedvehiclejourney.other)
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
