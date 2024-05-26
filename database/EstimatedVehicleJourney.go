package database

import (
	"database/sql"
	"fmt"
)

func InsertOrUpdateEstimatedVehicleJourney(db *sql.DB, values []interface{}) error {
	query := `
    INSERT INTO estimatedvehiclejourney (servicedelivery, recordedattime, lineref, directionref, datasource, datedvehiclejourneyref, vehiclemode, dataframeref, originref, destinationref, operatorref, vehicleref, cancellation, other)
    VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
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
        other = COALESCE(EXCLUDED.other, estimatedvehiclejourney.other);
    `

	stmt, err := db.Prepare(query)
	if err != nil {
		return fmt.Errorf("error preparing statement: %v", err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(values...)
	if err != nil {
		return fmt.Errorf("error executing statement: %v", err)
	}

	return nil
}
