package storage

import (
	"database/sql"
	"fmt"
	"net/url"

	_ "github.com/mattn/go-sqlite3"
)

func NewSQLite(path string) (*sql.DB, error) {
	params := url.Values{}

	params.Add("cache", "shared")
	params.Add("mode", "rwc")
	params.Add("_journal", "WAL")

	conn, err := sql.Open("sqlite3", fmt.Sprintf("file:%s?%s", path, params.Encode()))
	if err != nil {
		return nil, fmt.Errorf("failed to create conntection: %w", err)
	}

	// https://github.com/mattn/go-sqlite3/issues/209
	conn.SetMaxOpenConns(1)

	err = conn.Ping()
	if err != nil {
		return nil, fmt.Errorf("failed to ping connection: %w", err)
	}

	err = migrateSQLite(conn)
	if err != nil {
		return nil, fmt.Errorf("failed to run migrations: %w", err)
	}

	return conn, nil
}

func migrateSQLite(db *sql.DB) error {
	const schema = `
		CREATE TABLE IF NOT EXISTS state (
			region          TEXT       NOT NULL,
			log_group_name  TEXT       NOT NULL,
			log_stream_name TEXT       NOT NULL,
			next_token      TEXT       NOT NULL,
			updated_at      INTEGER    NOT NULL,

			PRIMARY KEY (region, log_group_name, log_stream_name)
		);
	`

	_, err := db.Exec(schema)
	if err != nil {
		return fmt.Errorf("failed to execute sql query: %w", err)
	}

	return nil
}
