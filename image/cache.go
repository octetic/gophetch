package image

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
	ccHeader := header.Get("Cache-Control")
	if ccHeader == "" {
		return Cache{
			Available:      false,
			MaxAge:         -1,
			Expires:        time.Time{},
			NoCache:        false,
			NoStore:        false,
			MustRevalidate: false,
		}, false
	}

	cache := Cache{}

	directives := strings.Split(ccHeader, ",")
	for _, directive := range directives {
		directive = strings.TrimSpace(directive)

		switch {
		case strings.HasPrefix(directive, "max-age"):
			value := strings.TrimPrefix(directive, "max-age=")
			age, err := strconv.Atoi(value)
			if err == nil {
				cache.MaxAge = age
				cache.Expires = time.Now().Add(time.Duration(age) * time.Second)
			}
		case directive == "no-cache":
			cache.NoCache = true
		case directive == "no-store":
			cache.NoStore = true
		case directive == "must-revalidate":
			cache.MustRevalidate = true
		}
	}

	cache.Available = true
	return cache, true
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
