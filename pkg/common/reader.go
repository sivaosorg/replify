package common

import (
	"bufio"
	"bytes"
	"io"
	"strings"
)

// ReadAll reads all data from an io.Reader and returns it as a single string.
//
// This function uses an `io.Copy` operation to efficiently read data from the provided
// `io.Reader` and write it to a `bytes.Buffer`. The resulting buffer content is then
// converted to a string and returned.
//
// Parameters:
//   - in: An `io.Reader` from which the function will read data. This can be any
//     type that implements the `io.Reader` interface, such as a file, standard input,
//     or a network connection.
//
// Returns:
//   - A string containing all the data read from the `io.Reader`.
//   - An error if any I/O operation fails during the copy process. If the input is
//     successfully read until EOF, the error returned is `nil`.
//
// Details:
//   - The function creates a `bytes.Buffer` to store the data read from the `io.Reader`.
//   - It uses the `io.Copy` function to transfer data from the `io.Reader` to the `bytes.Buffer`.
//     This approach is simple and efficient, leveraging built-in Go utilities for stream copying.
//   - After copying is complete, the data in the buffer is converted to a string and returned.
//
// Example Usage:
//
//	// Example: Reading from a file
//	file, err := os.Open("example.txt")
//	if err != nil {
//	    log.Fatal(err)
//	}
//	defer file.Close()
//
//	content, err := ReadAll(file)
//	if err != nil {
//	    log.Fatal(err)
//	}
//	fmt.Println(content)
//
//	// Example: Reading from standard input
//	fmt.Println("Enter some text (press Ctrl+D to end):")
//	content, err = ReadAll(os.Stdin)
//	if err != nil {
//	    log.Fatal(err)
//	}
//	fmt.Println("You entered:")
//	fmt.Println(content)
func ReadAll(in io.Reader) (string, error) {
	var buf bytes.Buffer
	if _, err := io.Copy(&buf, in); err != nil {
		return "", err
	}
	return buf.String(), nil
}

// SlurpAll reads all data from an io.Reader and returns it as a single string.
//
// This function reads raw bytes from the provided `io.Reader` in chunks and
// appends them to a buffer until EOF (end of file) is reached.
//
// Parameters:
//   - in: An `io.Reader` from which the function will read data. This could be any
//     type that implements the `io.Reader` interface, such as a file, standard input,
//     or a network connection.
//
// Returns:
//   - A string containing all the data read from the `io.Reader`.
//   - An error if any I/O operation fails (other than reaching EOF). If the input is
//     successfully read until EOF, the error returned is `nil`.
//
// Details:
//   - The function creates a buffer of size 1024 bytes to read chunks of data from
//     the `io.Reader`. This approach avoids reading the entire input into memory at once
//     and is suitable for handling large streams of data.
//   - The data read from each chunk is written to a `bytes.Buffer`, which efficiently
//     constructs the final string.
//   - If EOF is encountered during reading, the loop ends, and the accumulated data
//     in the buffer is returned as a string.
//   - Any other error during reading results in an immediate return of the error.
//
// Example Usage:
//
//	// Example: Reading from a file
//	file, err := os.Open("example.txt")
//	if err != nil {
//	    log.Fatal(err)
//	}
//	defer file.Close()
//
//	content, err := SlurpAll(file)
//	if err != nil {
//	    log.Fatal(err)
//	}
//	fmt.Println(content)
//
//	// Example: Reading from standard input
//	fmt.Println("Enter some text (press Ctrl+D to end):")
//	content, err = SlurpAll(os.Stdin)
//	if err != nil {
//	    log.Fatal(err)
//	}
//	fmt.Println("You entered:")
//	fmt.Println(content)
func SlurpAll(in io.Reader) (string, error) {
	buf := make([]byte, 1024)
	var out bytes.Buffer
	for {
		n, err := in.Read(buf)

		if err == io.EOF {
			break
		}
		if err != nil {
			return "", err
		}
		out.Write(buf[:n])
	}
	return out.String(), nil
}

// SlurpLines reads all lines from an io.Reader and returns them as a slice of strings.
//
// This function utilizes a buffered reader to read data line by line until EOF
// (end of file) is reached. Each line is appended to a slice of strings.
//
// Parameters:
//   - in: An `io.Reader` from which the function will read data. This can be any
//     type that implements the `io.Reader` interface, such as a file, standard input,
//     or a network connection.
//
// Returns:
//   - A slice of strings, where each element corresponds to a line read from the `io.Reader`.
//   - An error if any I/O operation fails (other than reaching EOF). If the input is
//     successfully read until EOF, the error returned is `nil`.
//
// Details:
//   - The function uses `bufio.NewReader` to create a buffered reader around the provided
//     `io.Reader`, allowing efficient reading of large input streams.
//   - It reads lines one by one using `ReadString('\n')`. Each line is appended to
//     a slice of strings. If a line does not end with a newline character before reaching EOF,
//     the remaining text is included as the last line in the slice.
//   - If an error occurs during reading (other than EOF), it stops reading and returns
//     the error along with the slice of lines read so far.
//
// Example Usage:
//
//	// Example: Reading from a file
//	file, err := os.Open("example.txt")
//	if err != nil {
//	    log.Fatal(err)
//	}
//	defer file.Close()
//
//	lines, err := SlurpLines(file)
//	if err != nil {
//	    log.Fatal(err)
//	}
//	for i, line := range lines {
//	    fmt.Printf("Line %d: %s", i+1, line)
//	}
//
//	// Example: Reading from standard input
//	fmt.Println("Enter some text (press Ctrl+D to end):")
//	lines, err = SlurpLines(os.Stdin)
//	if err != nil {
//	    log.Fatal(err)
//	}
//	fmt.Println("You entered:")
//	for _, line := range lines {
//	    fmt.Println(line)
//	}
func SlurpLines(in io.Reader) ([]string, error) {
	var lines []string
	scanner := bufio.NewReader(in)
	for {
		line, err := scanner.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			return lines, err
		}
		lines = append(lines, line)
	}
	return lines, nil
}

// SlurpLine reads all lines from an io.Reader and returns them as a single string.
//
// This function utilizes a buffered reader to read data line by line until EOF
// (end of file) is reached. It appends each line to a `strings.Builder`,
// which efficiently constructs the resulting string.
//
// Parameters:
//   - in: An `io.Reader` from which the function will read data. This could be any
//     type that implements the `io.Reader` interface, such as a file, standard input,
//     or a network connection.
//
// Returns:
//   - A string containing all the lines read from the `io.Reader`, concatenated together.
//   - An error if any I/O operation fails (other than reaching EOF). If the input is
//     successfully read until EOF, the error returned is `nil`.
//
// Details:
//   - The function uses `bufio.NewReader` to create a buffered reader around the provided
//     `io.Reader`. This allows efficient reading of large input streams.
//   - It reads lines one by one using `ReadString('\n')`. If a line does not end with a
//     newline character before reaching EOF, the remaining text is still included in the
//     result.
//   - If an error occurs during reading (other than EOF), it stops reading and returns
//     the error along with an empty string.
//
// Example Usage:
//
//	// Example: Reading from a file
//	file, err := os.Open("example.txt")
//	if err != nil {
//	    log.Fatal(err)
//	}
//	defer file.Close()
//
//	content, err := SlurpLine(file)
//	if err != nil {
//	    log.Fatal(err)
//	}
//	fmt.Println(content)
//
//	// Example: Reading from standard input
//	fmt.Println("Enter some text:")
//	content, err = SlurpLine(os.Stdin)
//	if err != nil {
//	    log.Fatal(err)
//	}
//	fmt.Println("You entered:", content)
//
// SlurpLine reads all lines from an io.Reader and returns them as a single string.
// It is used internally by ParseBufio.
func SlurpLine(in io.Reader) (string, error) {
	var lines strings.Builder
	scanner := bufio.NewReader(in)
	for {
		line, err := scanner.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			return "", err
		}
		lines.WriteString(line)
	}
	return lines.String(), nil
}
