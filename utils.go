package webthing

import (
	"time"
)

// Timestamp Get the current time.
//
// @return The current time in the form YYYY-mm-ddTHH:MM:SS+00.00
func Timestamp() string {
	now := time.Now().UTC().Format("2006-01-02T15:04:05")
	return now + "+00:00"
}
