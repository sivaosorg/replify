// Package strutil provides an extensive collection of string utility
// functions used throughout the replify library.
//
// The package operates exclusively on UTF-8 strings. Where character
// counting is relevant, functions work on Unicode code points (runes) rather
// than raw bytes, ensuring correct behaviour for multi-byte characters.
//
// # Emptiness Checks
//
//	strutil.IsEmpty(s)              // true when s is blank or whitespace-only
//	strutil.IsAnyEmpty(a, b, ...)   // true when any argument is empty
//	strutil.IsAllEmpty(a, b, ...)   // true when all arguments are empty
//	strutil.IsNotEmpty(s)           // convenience negation of IsEmpty
//
// # Trimming and Normalisation
//
// Functions cover left/right/both-side trimming, duplicate-whitespace
// collapse, Unicode normalisation, and removal of specific characters or
// substrings.
//
// # Case Conversion
//
// Beyond the standard ToUpper / ToLower helpers, the package provides
// ToCamelCase, ToSnakeCase, ToKebabCase, ToPascalCase, and ToTitleCase for
// identifier-style transformations.
//
// # Searching and Comparison
//
//	strutil.Contains(s, sub)        // substring test
//	strutil.ContainsAny(s, ...)     // any of the given substrings
//	strutil.StartsWith / EndsWith   // prefix / suffix checks
//	strutil.Count(s, sub)           // non-overlapping occurrence count
//	strutil.EqualFold(a, b)         // case-insensitive equality
//
// # Splitting and Joining
//
// Split, SplitN, and their trimming variants complement the standard library.
// JoinNonEmpty filters blank strings before joining, and Repeat wraps
// strings.Repeat with a rune-count guard.
//
// # Hashing
//
// Hash256 returns the SHA-256 hex digest of a string. The Len variable is an
// alias for utf8.RuneCountInString, providing a concise way to obtain the
// rune length of a string.
//
// # Integration
//
// strutil is a foundational dependency for several other packages in replify,
// including fj, match, encoding, and truncate.
//
// All functions in this package are pure and safe for concurrent use.
package strutil
