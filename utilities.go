package replify

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"

	"github.com/sivaosorg/replify/pkg/encoding"
	"github.com/sivaosorg/replify/pkg/randn"
)

// calculateSize calculates the size of the marshaled data.
// It uses encoding.Marshal to marshal the data and returns the length of the resulting byte slice.
// If an error occurs during marshaling, it returns 0.
func calculateSize(data any) int {
	_bytes, err := encoding.MarshalJSONb(data)
	if err != nil {
		return 0
	}
	return len(_bytes)
}

// compress compresses the given data using gzip and encodes it in base64.
// It first marshals the data using encoding.Marshal, then compresses the resulting byte slice
// using gzip. The compressed data is then encoded in base64 and returned as a string.
// If any error occurs during marshaling or compression, it returns an empty string.
func compress(data any) string {
	_bytes, err := encoding.MarshalJSONb(data)
	if err != nil {
		return ""
	}
	var buf bytes.Buffer
	gz := gzip.NewWriter(&buf)
	_, err = gz.Write(_bytes)
	if err != nil {
		return ""
	}
	err = gz.Close()
	if err != nil {
		return ""
	}
	return base64.StdEncoding.EncodeToString(buf.Bytes())
}

// decompress decompresses the given data using gzip and decodes it from base64.
// It first decodes the base64 encoded data using base64.StdEncoding.DecodeString,
// then decompresses the resulting byte slice using gzip. The decompressed data is
// then unmarshaled using encoding.Unmarshal and returned as an interface{}.
// If any error occurs during decoding or decompression, it returns nil.
func decompress(data string) any {
	_bytes, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return nil
	}
	gz, err := gzip.NewReader(bytes.NewReader(_bytes))
	if err != nil {
		return nil
	}
	defer gz.Close()
	var buf bytes.Buffer
	_, err = buf.ReadFrom(gz)
	if err != nil {
		return nil
	}
	var result any
	if err := encoding.UnmarshalBytes(buf.Bytes(), &result); err != nil {
		return nil
	}
	return result
}

// chunk takes a response represented as a map and returns a slice of byte slices,
// where each byte slice is a chunk of the JSON representation of the response.
// This is useful for streaming large responses in smaller segments.
// If the JSON encoding fails, it returns nil.
func chunk(data map[string]any) [][]byte {
	_bytes, err := encoding.MarshalJSONb(data)
	if err != nil {
		return nil
	}
	var chunks [][]byte
	for i := 0; i < len(_bytes); i += defaultChunkSize {
		end := i + defaultChunkSize
		if end > len(_bytes) {
			end = len(_bytes)
		}
		// Create a copy of the chunk to avoid referencing the underlying array.
		// This is important to ensure that each chunk is independent and can be
		// processed separately without affecting the others.
		chunk := make([]byte, end-i)
		copy(chunk, _bytes[i:end])
		chunks = append(chunks, chunk)
	}
	return chunks
}

// cryptoID generates a cryptographically secure random ID as a hexadecimal string.
// It uses 16 random bytes, which are then encoded to a hexadecimal string for easy representation.
//
// Returns:
//   - A string representing a secure random hexadecimal ID of 32 characters (since 16 bytes are used, and each byte
//     is represented by two hexadecimal characters).
//
// The function uses crypto/rand.Read to ensure cryptographic security in the generated ID, making it suitable for
// sensitive use cases such as API keys, session tokens, or any security-critical identifiers.
//
// Example:
//
//	id := cryptoID()
//	fmt.Println("Generated Crypto ID:", id)
//
// Notes:
//   - This function is suitable for use cases where high security is required in the generated ID.
//   - It is not recommended for use cases where deterministic or non-cryptographic IDs are preferred.
func cryptoID() string {
	return randn.CryptoID()
}

// jsonpass converts a Go value to its JSON string representation or returns the value directly if it is already a string.
//
// This function checks if the input data is a string; if so, it returns it directly.
// Otherwise, it marshals the input value `data` into a JSON string using the
// MarshalToString function. If an error occurs during marshalling, it returns an empty string.
//
// Parameters:
//   - `data`: The Go value to be converted to JSON, or a string to be returned directly.
//
// Returns:
//   - A string containing the JSON representation of the input value, or an empty string if an error occurs.
//
// Example:
//
//	jsonStr := jsonpass(myStruct)
func jsonpass(data any) string {
	return encoding.JSON(data)
}

// jsonpretty converts a Go value to its pretty-printed JSON string representation or returns the value directly if it is already a string.
//
// This function checks if the input data is a string; if so, it returns it directly.
// Otherwise, it marshals the input value `data` into a formatted JSON string using
// the MarshalIndent function. If an error occurs during marshalling, it returns an empty string.
//
// Parameters:
//   - `data`: The Go value to be converted to pretty-printed JSON, or a string to be returned directly.
//
// Returns:
//   - A string containing the pretty-printed JSON representation of the input value, or an empty string if an error occurs.
//
// Example:
//
//	jsonPrettyStr := jsonpretty(myStruct)
func jsonpretty(data any) string {
	return encoding.JSONPretty(data)
}
