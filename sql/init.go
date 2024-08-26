package sql

import (
	"database/sql"
	"fmt"
)

func InitDB(db *sql.DB) error {
	query := `
    CREATE TABLE IF NOT EXISTS contacts (
        id TEXT PRIMARY KEY,
        firstname TEXT,
        lastname TEXT,
        address TEXT,
        phone TEXT
    );`

	_, err := db.Exec(query)
	if err != nil {
		return fmt.Errorf("failed to create contacts table: %w", err)
	}

	_, err = db.Exec("CREATE INDEX IF NOT EXISTS idx_id ON contacts(id);")
	if err != nil {
		return fmt.Errorf("failed to create index on id: %w", err)
	}

	_, err = db.Exec("CREATE INDEX IF NOT EXISTS idx_name ON contacts(firstname);")
	if err != nil {
		return fmt.Errorf("failed to create index on name and phone: %w", err)
	}

	return nil
}
