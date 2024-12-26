package database

import (
	"database/sql"
	"errors"
	"math/rand"
)

func InsertOrUpdateEstimatedCall(db *sql.DB, estimatedValues []string) (int, string, error) {
	if len(estimatedValues) == 0 {
		return 0, "", errors.New("no estimated values provided")
	}

	id := rand.Intn(100)
	action := "tmp"
	return id, action, nil
}
