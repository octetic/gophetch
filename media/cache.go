package media

import (
	"net/http"
	"strconv"
	"strings"
	"time"
)

// Cache holds the information extracted from Cache-Control and Expires headers
type Cache struct {
	Available      bool
	MaxAge         int
	Expires        time.Time
	NoCache        bool
	NoStore        bool
	MustRevalidate bool
}

// ParseCacheHeader takes a http.Header and returns the parsed Cache
func ParseCacheHeader(header http.Header) Cache {
	cache, cacheFound := parseCacheControlHeader(header)

	var expiresFound bool
	if cache.MaxAge < 0 {
		cache.Expires, expiresFound = parseExpiresHeader(header)
	}

	cache.Available = cacheFound || expiresFound
	return cache
}

func parseCacheControlHeader(header http.Header) (Cache, bool) {
	ccInfo := Cache{
		Available:      false,
		MaxAge:         -1,
		Expires:        time.Time{},
		NoCache:        false,
		NoStore:        false,
		MustRevalidate: false,
	}

	// Extract Cache-Control header
	cacheControl := header.Get("Cache-Control")
	if cacheControl == "" {
		return ccInfo, false
	}

	for _, directive := range strings.Split(cacheControl, ",") {
		parts := strings.SplitN(strings.TrimSpace(directive), "=", 2)
		switch parts[0] {
		case "max-age":
			ccInfo.MaxAge = parseMaxAge(parts[1])
		case "no-cache":
			ccInfo.NoCache = true
		case "no-store":
			ccInfo.NoStore = true
		case "must-revalidate":
			ccInfo.MustRevalidate = true
		}
	}

	ccInfo.Available = true
	return ccInfo, true
}

// parseMaxAge parses the Max-Age directive from Cache-Control header.
func parseMaxAge(value string) int {
	maxAge, err := strconv.Atoi(value)
	if err != nil {
		return -1
	}
	return maxAge
}

// parseExpiresHeader takes an http.Header and returns the parsed Expires time
func parseExpiresHeader(header http.Header) (time.Time, bool) {
	expiresHeader := header.Get("Expires")
	if expiresHeader == "" {
		return time.Time{}, false
	}

	expires, err := http.ParseTime(expiresHeader)
	if err != nil {
		return time.Time{}, false
	}

	return expires, true
}
