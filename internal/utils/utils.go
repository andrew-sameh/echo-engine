package utils

import (
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

// Now returns the current time as a pgtype.Timestamp UTC.
func PgTimeNow() pgtype.Timestamp {
	return pgtype.Timestamp{Time: time.Now().UTC(), Valid: true}
}

func PgTimeNowLocal() pgtype.Timestamp {
	return pgtype.Timestamp{Time: time.Now().Local(), Valid: true}
}
