// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: query.sql

package db

import (
	"context"
	"time"
)

const pingDB = `-- name: PingDB :one
SELECT NOW()
`

func (q *Queries) PingDB(ctx context.Context) (time.Time, error) {
	row := q.db.QueryRowContext(ctx, pingDB)
	var now time.Time
	err := row.Scan(&now)
	return now, err
}
