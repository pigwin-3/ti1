package database

import (
	"context"
	"database/sql"
	"fmt"
	"sync"
)

type EVJResult struct {
	ID     int
	Action string
	Error  error
	Index  int // To maintain order
}

// PreparedStatements holds reusable prepared statements
type PreparedStatements struct {
	evjStmt *sql.Stmt
	ecStmt  *sql.Stmt
	rcStmt  *sql.Stmt
	mu      sync.Mutex
}

func NewPreparedStatements(db *sql.DB) (*PreparedStatements, error) {
	evjQuery := `
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

	evjStmt, err := db.Prepare(evjQuery)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare EVJ statement: %w", err)
	}

	return &PreparedStatements{
		evjStmt: evjStmt,
	}, nil
}

func (ps *PreparedStatements) Close() {
	if ps.evjStmt != nil {
		ps.evjStmt.Close()
	}
	if ps.ecStmt != nil {
		ps.ecStmt.Close()
	}
	if ps.rcStmt != nil {
		ps.rcStmt.Close()
	}
}

func InsertOrUpdateEstimatedVehicleJourney(db *sql.DB, values []interface{}) (int, string, error) {
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

	var action string
	var id int
	err := db.QueryRow(query, values...).Scan(&action, &id)
	if err != nil {
		return 0, "", fmt.Errorf("error executing EVJ statement: %w", err)
	}

	return id, action, nil
}

// BatchInsertEVJ processes multiple EVJ inserts concurrently
func BatchInsertEVJ(ctx context.Context, db *sql.DB, batch [][]interface{}, workerCount int) ([]EVJResult, error) {
	if len(batch) == 0 {
		return nil, nil
	}

	results := make([]EVJResult, len(batch))
	jobs := make(chan int, len(batch))
	var wg sync.WaitGroup

	// Start workers
	for w := 0; w < workerCount; w++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for idx := range jobs {
				select {
				case <-ctx.Done():
					return
				default:
					id, action, err := InsertOrUpdateEstimatedVehicleJourney(db, batch[idx])
					results[idx] = EVJResult{
						ID:     id,
						Action: action,
						Error:  err,
						Index:  idx,
					}
				}
			}
		}()
	}

	// Send jobs
	for i := range batch {
		jobs <- i
	}
	close(jobs)

	wg.Wait()
	return results, nil
}
