package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"
)

type State struct {
	db *sql.DB
}

func NewState(db *sql.DB) *State {
	return &State{
		db: db,
	}
}

func (s *State) GetNextToken(ctx context.Context, region, logGroupName, logStreamName string) (string, error) {
	const query = `SELECT next_token FROM state WHERE region = ? and log_group_name = ? and log_stream_name = ?`

	var nextToken string

	err := s.db.QueryRowContext(ctx, query, region, logGroupName, logStreamName).Scan(&nextToken)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return "", fmt.Errorf("failed to execute sql query: %w", err)
	}

	return nextToken, nil
}

func (s *State) SetNextToken(ctx context.Context, region, logGroupName, logStreamName, nextToken string) error {
	const query = `
		INSERT INTO state (region, log_group_name, log_stream_name, next_token, updated_at)
		VALUES (?, ?, ?, ?, ?)
		ON CONFLICT (region, log_group_name, log_stream_name) DO UPDATE SET
			next_token = excluded.next_token,
			updated_at = excluded.updated_at
	`

	_, err := s.db.ExecContext(ctx, query, region, logGroupName, logStreamName, nextToken, time.Now().UTC().Unix())
	if err != nil {
		return fmt.Errorf("failed to execute sql query: %w", err)
	}

	return nil
}
