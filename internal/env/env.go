package env

import (
	"os"
	"strconv"
	"strings"
	"time"
)

// String returns the env variable value trimmed from whitespace and optional surrounding quotes.
func String(key, def string) string {
	v := strings.TrimSpace(os.Getenv(key))
	if v == "" {
		return def
	}
	return strings.Trim(v, "\"")
}

func Int(key string, def int) int {
	v := strings.TrimSpace(os.Getenv(key))
	v = strings.Trim(v, "\"")
	if v == "" {
		return def
	}
	i, err := strconv.Atoi(v)
	if err != nil {
		return def
	}
	return i
}

// Bool parse boolean string to bool
// True: "1", "t", "T", "true", "TRUE", "True"
// False: "0", "f", "F", "false", "FALSE", "False"
func Bool(key string, def bool) bool {
	v := strings.TrimSpace(os.Getenv(key))
	v = strings.Trim(v, "\"")
	if v == "" {
		return def
	}
	b, err := strconv.ParseBool(v)
	if err != nil {
		return def
	}
	return b
}

// Duration reads a time.Duration env var (Go duration format).
func Duration(key string, def time.Duration) time.Duration {
	v := strings.TrimSpace(os.Getenv(key))
	v = strings.Trim(v, "\"")
	if v == "" {
		return def
	}
	d, err := time.ParseDuration(v)
	if err != nil {
		return def
	}
	return d
}
