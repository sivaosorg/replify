// Package match provides wildcard glob pattern matching for strings.
//
// Two wildcard characters are supported:
//
//	*   matches any sequence of characters, including the empty sequence.
//	?   matches exactly one arbitrary character.
//
// A backslash escapes the character that follows it, allowing literal '*'
// and '?' to appear in patterns. All matching is done on Unicode code points
// so multi-byte UTF-8 sequences are handled correctly.
//
// # Basic Matching
//
//	match.Match("hello", "h*o")      // true
//	match.Match("hello", "h?llo")    // true
//	match.Match("hello", "*")        // true  (star matches everything)
//	match.Match("hello", "world")    // false
//
// # Complexity-Limited Matching
//
// MatchLimit guards against adversarial inputs by capping the number of
// recursive wildcard expansions. It returns both a match result and a
// stopped flag that is set when the complexity budget is exhausted:
//
//	matched, stopped := match.MatchLimit(str, pattern, 100)
//	if stopped {
//	    // pattern was too complex; treat as no-match
//	}
//
// # Range Queries
//
// WildcardPatternLimits derives the lexicographically smallest and largest
// strings that could match a given pattern. This is useful for efficiently
// narrowing a sorted index before applying a full wildcard match:
//
//	min, max := match.WildcardPatternLimits("user:*:profile")
//
// match is used internally by the fj package to evaluate path conditions
// in JSON queries and is exposed as a standalone package for use in other
// filtering contexts within replify.
//
// All functions are stateless and safe for concurrent use by multiple
// goroutines.
package match
