package webthing

import (
	"errors"
	"net/http"
	"regexp"
	"time"
)

// Timestamp Get the current time.
//
// @return The current time in the form YYYY-mm-ddTHH:MM:SS+00.00
func Timestamp() string {
	now := time.Now().UTC().Format("2006-01-02T15:04:05")
	return now + "+00:00"
}

func trimSlash(path string) string {
	l := len(path)
	if l != 1 && path[l-1:] == "/" {
		return path[:l-1]
	}
	return path
}

func resource(path string) (string, error) {
	m := validPath().FindStringSubmatch(path)
	if m == nil {
		return "", errors.New(" Invalid path! ")
	}
	return m[2], nil // The resource is the second subexpression.
}

func validPath() *regexp.Regexp {
	// return regexp.MustCompile(`\/(properties|actions|events)\/([a-zA-Z0-9]+)$`)
	return regexp.MustCompile(`(properties|actions|events)\/([a-zA-Z0-9]+)`)
}

// corsResponse Add necessary CORS headers to response.
//
// @param response Response to add headers to
// @return The Response object.
func corsResponse(w http.ResponseWriter) http.ResponseWriter {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept")
	w.Header().Set("Access-Control-Allow-Methods", "GET, HEAD, PUT, POST, DELETE")
	return w
}

//jsonResponse Add json headers to response.
func jsonResponse(w http.ResponseWriter) http.ResponseWriter {
	w.Header().Set("Content-Type", "application/json")
	return w
}
