package codegen

const (
	// CharsetNumeric contains only the digits 0-9.
	CharsetNumeric Charset = "0123456789"

	// CharsetAlphaLower contains only lowercase letters a-z.
	CharsetAlphaLower Charset = "abcdefghijklmnopqrstuvwxyz"

	// CharsetAlphaUpper contains only uppercase letters A-Z.
	CharsetAlphaUpper Charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"

	// CharsetAlpha contains both lowercase and uppercase letters (a-z, A-Z).
	CharsetAlpha Charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

	// CharsetAlphanumeric contains digits, lowercase letters, and uppercase letters.
	// This is the default character set.
	CharsetAlphanumeric Charset = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

	// CharsetAlphanumericUpper contains digits and uppercase letters only.
	// It is commonly used for easy-to-read order codes.
	CharsetAlphanumericUpper Charset = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ"

	// CharsetAlphanumericLower contains digits and lowercase letters only.
	CharsetAlphanumericLower Charset = "0123456789abcdefghijklmnopqrstuvwxyz"
)
