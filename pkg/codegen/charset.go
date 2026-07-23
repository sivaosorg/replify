package codegen

// String returns the string representation of the Charset.
func (c Charset) String() string {
	return string(c)
}

// Len returns the number of characters in the Charset.
func (c Charset) Len() int {
	return len([]rune(string(c)))
}
