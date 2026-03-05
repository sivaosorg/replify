package sysx

import (
	"os"
	"strconv"
	"strings"
)

// ///////////////////////////
// Section: Basic env access
// ///////////////////////////

// GetEnv returns the value of the environment variable named by key.
//
// If the variable is not set, or is set to an empty string, GetEnv returns
// the provided fallback value instead.
//
// Parameters:
//   - `key`:      the name of the environment variable.
//   - `fallback`: the value to return when the variable is absent or empty.
//
// Returns:
//
//	A string: either the current value of the variable or fallback.
//
// Example:
//
//	port := sysx.GetEnv("PORT", "8080")
func GetEnv(key, fallback string) string {
	if v, ok := os.LookupEnv(key); ok && v != "" {
		return v
	}
	return fallback
}

// MustGetEnv returns the value of the environment variable named by key and
// panics if the variable is not set or is empty.
//
// Use this function only when the absence of the variable is considered an
// unrecoverable configuration error (e.g. during application bootstrap).
//
// Parameters:
//   - `key`: the name of the environment variable.
//
// Returns:
//
//	A non-empty string containing the value of the variable.
//
// Example:
//
//	dsn := sysx.MustGetEnv("DATABASE_URL")
func MustGetEnv(key string) string {
	v, ok := os.LookupEnv(key)
	if !ok || v == "" {
		panic("sysx: required environment variable not set or empty: " + key)
	}
	return v
}

// HasEnv reports whether the environment variable named by key exists and
// is non-empty.
//
// A variable that is explicitly set to an empty string is treated as absent
// for the purposes of this function.
//
// Parameters:
//   - `key`: the name of the environment variable.
//
// Returns:
//
//	A boolean value:
//	 - true  when the variable is set and its value is non-empty;
//	 - false otherwise.
//
// Example:
//
//	if sysx.HasEnv("DEBUG") {
//	    fmt.Println("debug mode enabled")
//	}
func HasEnv(key string) bool {
	v, ok := os.LookupEnv(key)
	return ok && v != ""
}

// ///////////////////////////
// Section: Env mutation
// ///////////////////////////

// SetEnv sets the environment variable named by key to value.
//
// It delegates directly to os.Setenv and propagates any error.
//
// Parameters:
//   - `key`:   the name of the environment variable to set.
//   - `value`: the value to assign.
//
// Returns:
//
//	An error if the variable could not be set, or nil on success.
//
// Example:
//
//	if err := sysx.SetEnv("LOG_LEVEL", "debug"); err != nil {
//	    log.Fatal(err)
//	}
func SetEnv(key, value string) error {
	return os.Setenv(key, value)
}

// UnsetEnv removes the environment variable named by key from the process
// environment.
//
// It delegates directly to os.Unsetenv and propagates any error.
//
// Parameters:
//   - `key`: the name of the environment variable to remove.
//
// Returns:
//
//	An error if the variable could not be removed, or nil on success.
//
// Example:
//
//	if err := sysx.UnsetEnv("TEMP_TOKEN"); err != nil {
//	    log.Fatal(err)
//	}
func UnsetEnv(key string) error {
	return os.Unsetenv(key)
}

// ///////////////////////////
// Section: Typed env helpers
// ///////////////////////////

// GetEnvInt returns the value of the environment variable named by key parsed
// as a base-10 integer.
//
// If the variable is not set, is empty, or cannot be parsed as an integer,
// the provided fallback value is returned.
//
// Parameters:
//   - `key`:      the name of the environment variable.
//   - `fallback`: the integer to return when the variable is absent, empty, or non-numeric.
//
// Returns:
//
//	An int: either the parsed variable value or fallback.
//
// Example:
//
//	workers := sysx.GetEnvInt("WORKERS", 4)
func GetEnvInt(key string, fallback int) int {
	v, ok := os.LookupEnv(key)
	if !ok || v == "" {
		return fallback
	}
	n, err := strconv.Atoi(v)
	if err != nil {
		return fallback
	}
	return n
}

// GetEnvBool returns the value of the environment variable named by key parsed
// as a boolean.
//
// The following string values (case-insensitive) are recognised:
//   - true  strings: "1", "true", "yes", "on"
//   - false strings: "0", "false", "no", "off"
//
// If the variable is not set, is empty, or contains an unrecognised value,
// the provided fallback value is returned.
//
// Parameters:
//   - `key`:      the name of the environment variable.
//   - `fallback`: the boolean to return when the variable cannot be resolved.
//
// Returns:
//
//	A bool: either the parsed variable value or fallback.
//
// Example:
//
//	debug := sysx.GetEnvBool("DEBUG", false)
func GetEnvBool(key string, fallback bool) bool {
	v, ok := os.LookupEnv(key)
	if !ok || v == "" {
		return fallback
	}
	result, ok := parseBoolString(strings.ToLower(strings.TrimSpace(v)))
	if !ok {
		return fallback
	}
	return result
}

// GetEnvSlice returns the value of the environment variable named by key
// split by the provided separator.
//
// If the variable is not set or is empty, a nil slice is returned. Empty
// elements produced by consecutive separators are preserved in the result.
//
// Parameters:
//   - `key`: the name of the environment variable.
//   - `sep`: the separator string used to split the variable's value.
//
// Returns:
//
//	A []string containing the split elements, or nil when the variable is absent or empty.
//
// Example:
//
//	hosts := sysx.GetEnvSlice("HOSTS", ",") // ["a.example.com", "b.example.com"]
func GetEnvSlice(key, sep string) []string {
	v, ok := os.LookupEnv(key)
	if !ok || v == "" {
		return nil
	}
	return strings.Split(v, sep)
}

// ///////////////////////////
// Section: Env introspection
// ///////////////////////////

// Environ returns a copy of all environment variables in the form
// "key=value".
//
// It delegates directly to os.Environ. The returned slice is a snapshot;
// subsequent changes to the process environment are not reflected.
//
// Returns:
//
//	A []string where each element is a "KEY=value" pair.
//
// Example:
//
//	for _, e := range sysx.Environ() {
//	    fmt.Println(e)
//	}
func Environ() []string {
	return os.Environ()
}

// EnvMap returns all current environment variables as a map from key to
// value.
//
// Each entry in os.Environ is split on the first '=' character. Variables
// without a value portion are stored with an empty string value. The
// returned map is a snapshot; subsequent changes to the process environment
// are not reflected.
//
// Returns:
//
//	A non-nil map[string]string containing all current environment variables.
//
// Example:
//
//	env := sysx.EnvMap()
//	fmt.Println(env["HOME"])
func EnvMap() map[string]string {
	pairs := os.Environ()
	m := make(map[string]string, len(pairs))
	for _, pair := range pairs {
		idx := strings.IndexByte(pair, '=')
		if idx < 0 {
			m[pair] = ""
		} else {
			m[pair[:idx]] = pair[idx+1:]
		}
	}
	return m
}
