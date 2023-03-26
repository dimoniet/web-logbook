package driver

import (
	"context"
	"database/sql"
	"encoding/json"
	"time"

	_ "embed"

	"github.com/vsimakhin/web-logbook/internal/models"
	_ "modernc.org/sqlite"
)

func OpenDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("sqlite", dsn)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	err = validateDB(db)
	if err != nil {
		return nil, err
	}

	return db, nil
}

//go:embed db.structure
var structure string

// validateDB creates db structure in case it's a first run and the schema is empty
func validateDB(db *sql.DB) error {
	var err error
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// run sql from db.structure
	_, err = db.ExecContext(ctx, structure)
	if err != nil {
		return err
	}

	// check settings table it's not empty
	err = checkSettingsTable(db)
	if err != nil {
		return err
	}

	return nil
}

// checkSettingsTable verifies the proper transition to the new settings table
func checkSettingsTable(db *sql.DB) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var s models.Settings

	// check settings table (first app run)
	var rowsCount int
	row := db.QueryRowContext(ctx, "SELECT COUNT(*) FROM settings2")

	err := row.Scan(&rowsCount)
	if err != nil {
		return err
	}

	if rowsCount == 0 {
		// default values
		s.OwnerName = "Owner Name"
		s.SignatureText = "I certify that the entries in this log are true."

		out, err := json.Marshal(s)
		if err != nil {
			return err
		}

		_, err = db.ExecContext(ctx, "INSERT INTO settings2 (id, settings) VALUES (0, ?)", string(out))
		if err != nil {
			return err
		}
	}

	return nil
}
