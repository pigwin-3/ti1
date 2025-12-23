package database

import (
	"context"
	"crypto/md5"
	"database/sql"
	"encoding/hex"
	"fmt"
	"sync"
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

	// Get the MD5 hash from Valkey
	retrievedHash, err := valki.GetValkeyValue(ctx, valkeyClient, key)
	if err != nil {
		return 0, "", fmt.Errorf("failed to get value from Valkey: %w", err)
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

		err = valki.SetValkeyValue(ctx, valkeyClient, key, hashString)
		if err != nil {
			return 0, "", fmt.Errorf("failed to set value in Valkey: %w", err)
		}

		var action string
		var id int
		err = db.QueryRowContext(ctx, query, values...).Scan(&action, &id)
		if err != nil {
			return 0, "", fmt.Errorf("error executing statement: %w", err)
		}
		return id, action, nil
	}
	return 0, "none", nil
}

// BatchInsertRecordedCalls processes multiple recorded calls concurrently
func BatchInsertRecordedCalls(ctx context.Context, db *sql.DB, batch [][]interface{}, valkeyClient valkey.Client, workerCount int) ([]CallResult, error) {
	if len(batch) == 0 {
		return nil, nil
	}

	results := make([]CallResult, len(batch))
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
					id, action, err := InsertOrUpdateRecordedCall(ctx, db, batch[idx], valkeyClient)
					results[idx] = CallResult{
						ID:     id,
						Action: action,
						Error:  err,
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
