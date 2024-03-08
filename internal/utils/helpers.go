package utils

import (
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

func UUIDToString(uuid pgtype.UUID) string {
	s := fmt.Sprintf(
		"%x-%x-%x-%x-%x",
		uuid.Bytes[0:4],
		uuid.Bytes[4:6],
		uuid.Bytes[6:8],
		uuid.Bytes[8:10],
		uuid.Bytes[10:16],
	)

	return s
}

func StringToTime(timestamp string, returnNow bool) time.Time {
	t, err := time.Parse(time.RFC3339, timestamp)
	if err != nil {
		if returnNow {
			return time.Now().UTC()
		}
		
		return time.Time{}
	}

	return t
}
func StringToPgTime(timestamp string, returnNow bool) pgtype.Timestamp {
	t, err := time.Parse(time.RFC3339, timestamp)
	if err != nil {
		if returnNow {
			return pgtype.Timestamp{Time: time.Now().UTC(), Valid: true}
		}
		
		return pgtype.Timestamp{Valid: false}
	}

	return pgtype.Timestamp{Time: t, Valid: true}
}
