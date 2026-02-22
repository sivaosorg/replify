package common

import (
	"bytes"
	"io"
)

// TeeCopy reads data from an io.Reader and writes the data simultaneously to an io.Writer,
// while also returning the logged data as a string.
//
// This function utilizes `io.TeeReader` to duplicate the data read from the provided `io.Reader`.
// A copy of the data is written to the provided `io.Writer`, while another copy is logged
// into a `bytes.Buffer`.
//
// Parameters:
//   - in: An `io.Reader` from which data will be read.
//   - out: An `io.Writer` to which the read data will be written.
//
// Returns:
//   - A string containing all the data read from the `io.Reader` as logged in the `bytes.Buffer`.
//   - An error if any I/O operation fails.
//
// Details:
//   - The function creates a `bytes.Buffer` to log a copy of the data.
//   - It uses `io.TeeReader` to split the data read from `in` between logging to the buffer
//     and passing the data to the consumer (`out`).
//   - A loop is used to read data in chunks and write it to the `io.Writer`.
//   - The function ensures that both reading and writing are handled efficiently.
//
// Example Usage:
//
//	// Example: Reading from a file and writing to another file while logging
//	inFile, err := os.Open("input.txt")
//	if err != nil {
//	    log.Fatal(err)
//	}
//	defer inFile.Close()
//
//	outFile, err := os.Create("output.txt")
//	if err != nil {
//	    log.Fatal(err)
//	}
//	defer outFile.Close()
//
//	loggedData, err := TeeCopy(inFile, outFile)
//	if err != nil {
//	    log.Fatal(err)
//	}
//	fmt.Println("Logged Data:")
//	fmt.Println(loggedData)
//
//	// Example: Reading from standard input and writing to standard output
//	fmt.Println("Enter some text (press Ctrl+D to end):")
//	loggedData, err = TeeCopy(os.Stdin, os.Stdout)
//	if err != nil {
//	    log.Fatal(err)
//	}
//	fmt.Println("\nLogged Data:")
//	fmt.Println(loggedData)
func TeeCopy(in io.Reader, out io.Writer) (string, error) {
	var lines bytes.Buffer
	reader := io.TeeReader(in, &lines)
	buffer := make([]byte, 256)
	for {
		n, err := reader.Read(buffer)
		if n > 0 {
			if _, writeErr := out.Write(buffer[:n]); writeErr != nil {
				return "", writeErr
			}
		}
		if err == io.EOF {
			break
		}
		if err != nil {
			return "", err
		}
	}
	return lines.String(), nil
}

// TeeTap reads data from an io.Reader and logs the data simultaneously, returning the logged data as a string.
//
// This function utilizes `io.TeeReader` to split the data read from the provided `io.Reader`.
// The `io.TeeReader` writes a copy of all read data into a `bytes.Buffer` while continuing to pass
// the data through as it is read. The function consumes all data from the `io.Reader` but discards it,
// keeping only the logged copy.
//
// Parameters:
//   - in: An `io.Reader` from which the function will read data. This can be any type that implements
//     the `io.Reader` interface, such as a file, standard input, or a network connection.
//
// Returns:
//   - A string containing all the data read from the `io.Reader` as logged in the `bytes.Buffer`.
//   - An error if any I/O operation fails (other than reaching EOF). If the input is successfully
//     read until EOF, the error returned is `nil`.
//
// Details:
//   - The function creates a `bytes.Buffer` to store a copy of the data being read from the `io.Reader`.
//   - An `io.TeeReader` is used to duplicate the data read: one copy is written to the buffer (`logs`),
//     and the other is read into a temporary buffer to consume the input stream.
//   - The function uses a loop to read from the `io.TeeReader` in chunks of 256 bytes until EOF is reached.
//   - The accumulated data in the `bytes.Buffer` is returned as a string after the reading is complete.
//
// Example Usage:
//
//	// Example: Reading from a file and logging data
//	file, err := os.Open("example.txt")
//	if err != nil {
//	    log.Fatal(err)
//	}
//	defer file.Close()
//
//	loggedData, err := TeeTap(file)
//	if err != nil {
//	    log.Fatal(err)
//	}
//	fmt.Println("Logged Data:")
//	fmt.Println(loggedData)
//
//	// Example: Reading from standard input
//	fmt.Println("Enter some text (press Ctrl+D to end):")
//	loggedData, err = TeeTap(os.Stdin)
//	if err != nil {
//	    log.Fatal(err)
//	}
//	fmt.Println("You entered (logged):")
//	fmt.Println(loggedData)
func TeeTap(in io.Reader) (string, error) {
	var lines bytes.Buffer
	reader := io.TeeReader(in, &lines)
	buffer := make([]byte, 256)
	for {
		_, err := reader.Read(buffer)
		if err == io.EOF {
			break
		}
		if err != nil {
			return "", err
		}
	}
	return lines.String(), nil
}
