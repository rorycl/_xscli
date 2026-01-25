package db

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "modernc.org/sqlite" // Import the pure Go SQLite driver
)

// Presently not sure whether GetBankTransactionUpdatedTime is needed.

// GetBankTransactionUpdatedTime retrieves the last updated timestamp for a single bank transaction.
func (db *DB) GetBankTransactionUpdatedTime(uuid string) (time.Time, error) {
	var updatedTime time.Time
	query := `SELECT updated_at FROM bank_transactions WHERE id = ?;`
	err := db.QueryRowContext(context.Background(), query, uuid).Scan(&updatedTime)
	if err != nil {
		if err == sql.ErrNoRows {
			return time.Time{}, fmt.Errorf("record not found in local database")
		}
		return time.Time{}, err
	}
	return updatedTime, nil
}
