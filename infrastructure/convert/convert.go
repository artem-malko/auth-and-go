package convert

import (
	"time"

	"gopkg.in/guregu/null.v3"
)

// NewStringPointer return pointer to string, if it is not empty
func NewStringPointer(s string) *string {
	if s == "" {
		return nil
	}

	return &s
}

// NewStringPointer return pointer to time, if it is not zero time
func NewTimePointer(t time.Time) *time.Time {
	if t.IsZero() {
		return nil
	}

	return &t
}

// NewSQLNullString return sql.NullString with value, if string is not empty
func NewSQLNullString(s string) null.String {
	return null.NewString(s, s != "")
}

// NewSQLNullTime return sql.NewSQLNullTime with value, if time is not zero
func NewSQLNullTime(t time.Time) null.Time {
	return null.NewTime(t, !t.IsZero())
}
