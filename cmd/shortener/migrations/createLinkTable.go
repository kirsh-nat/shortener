package migrations

import (
	"database/sql"
	"log"
)

func CreateLinkTable(db *sql.DB) {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS links (
			id SERIAL PRIMARY KEY,
			original_url TEXT NOT NULL,
			short_url TEXT NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			user_id TEXT DEFAULT NULL,
			deleted BOOLEAN DEFAULT FALSE
		);	
		CREATE UNIQUE INDEX IF NOT EXISTS short_url_unique ON links (short_url);
	`)

	if err != nil {
		log.Fatal(err)
	}
}
