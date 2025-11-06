package common

import (
	"net/http"
	"strconv"
)

// ParseIntParam retrieves an integer query parameter from the HTTP request.
func ParseIntParam(r *http.Request, param string, defaultValue int) int {
	valueStr := r.URL.Query().Get(param)
	if valueStr == "" {
		return defaultValue
	}

	value, err := strconv.Atoi(valueStr)
	if err != nil {
		return defaultValue
	}

	return value
}

// ParseOffsetLimit parses offset and limit query parameters with validation
func ParseOffsetLimit(r *http.Request) (int, int) {
	offset := ParseIntParam(r, "offset", 0)
	limit := ParseIntParam(r, "limit", 10)

	// Validate and constrain limit
	if limit < 1 {
		limit = 1
	}
	if limit > 100 {
		limit = 100
	}

	return offset, limit
}
