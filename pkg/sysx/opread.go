package sysx

import (
	"bufio"
	"bytes"
	"io"
	"os"
)

// ReadBytes reads the entire contents of the file at path and returns them as
// a byte slice.
//
// Parameters:
//   - `path`: the file system path to read.
//
// Returns:
//
//	([]byte, error): the file contents and nil on success, or nil and a
//	non-nil error if the file does not exist or cannot be read.
//
// Example:
//
//	data, err := sysx.ReadBytes("/etc/hosts")
//	if err != nil {
//	    log.Fatal(err)
//	}
//	fmt.Printf("%d bytes read\n", len(data))
func ReadBytes(path string) ([]byte, error) {
	return os.ReadFile(path)
}

// ReadString reads the entire contents of the file at path and returns
// them as a string.
//
// Parameters:
//   - `path`: the file system path to read.
//
// Returns:
//
//	(string, error): the file contents and nil on success, or an empty string
//	and a non-nil error if the file does not exist or cannot be read.
//
// Example:
//
//	content, err := sysx.ReadString("/etc/hostname")
//	if err != nil {
//	    log.Fatal(err)
//	}
//	fmt.Println(strings.TrimSpace(content))
func ReadString(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// ReadLines reads the file at path and returns its contents as a slice of
// strings, one element per line. Line endings ("\n" and "\r\n") are stripped
// from each element. An empty file returns a non-nil empty slice.
//
// The file is read with a buffered scanner, making it efficient even for
// large files as long as no single line exceeds bufio.MaxScanTokenSize
// (64 KiB by default).
//
// Parameters:
//   - `path`: the file system path to read.
//
// Returns:
//
//	([]string, error): lines of the file and nil on success, or a partial
//	result and a non-nil error on failure.
//
// Example:
//
//	lines, err := sysx.ReadLines("/var/log/app.log")
//	for _, l := range lines {
//	    fmt.Println(l)
//	}
func ReadLines(path string) ([]string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	lines := make([]string, 0, 64)
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return lines, err
	}
	return lines, nil
}

// StreamLines opens the file at path and calls handler for each line in order.
// Processing stops immediately and the handler's error is returned when handler
// returns a non-nil value. A bufio.Scanner error encountered during reading is
// also returned. Line endings are stripped before handler is invoked.
//
// StreamLines is designed for memory-efficient processing of large files: only
// one line is held in memory at a time.
//
// Parameters:
//   - `path`:    the file system path to read.
//   - `handler`: the function called for each line; return a non-nil error to stop.
//
// Returns:
//
//	An error if the file could not be opened, the scanner failed, or handler
//	returned a non-nil error; nil on success.
//
// Example:
//
//	count := 0
//	err := sysx.StreamLines("/var/log/access.log", func(line string) error {
//	    if strings.Contains(line, "ERROR") {
//	        count++
//	    }
//	    return nil
//	})
func StreamLines(path string, handler func(string) error) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		if err := handler(scanner.Text()); err != nil {
			return err
		}
	}
	return scanner.Err()
}

// CountLines counts the number of lines in the file at path.
//
// Counting is performed by reading the file in 32 KiB chunks and counting
// newline characters ("\n"), making it significantly faster than line-by-line
// scanning and avoiding line length limits. An empty file returns 0.
//
// Parameters:
//   - `path`: the file system path to count lines in.
//
// Returns:
//
//	(int, error): the number of lines and nil on success, or 0 and a non-nil
//	error if the file cannot be opened or read.
//
// Example:
//
//	n, err := sysx.CountLines("/var/log/app.log")
//	if err != nil {
//	    log.Fatal(err)
//	}
//	fmt.Printf("%d lines\n", n)
func CountLines(path string) (int, error) {
	f, err := os.Open(path)
	if err != nil {
		return 0, err
	}
	defer f.Close()

	count := 0
	buf := make([]byte, 32768)
	lineSep := []byte{'\n'}

	for {
		n, err := f.Read(buf)
		count += bytes.Count(buf[:n], lineSep)
		if err == io.EOF {
			break
		}
		if err != nil {
			return count, err
		}
	}
	return count, nil
}

// Head reads the first n lines of the file at path and returns them as a
// slice of strings with line endings stripped. If the file has fewer than n
// lines, all lines are returned. n must be non-negative; passing 0 returns an
// empty (non-nil) slice with no error.
//
// Parameters:
//   - `path`: the file system path to read.
//   - `n`:    the maximum number of lines to return.
//
// Returns:
//
//	([]string, error): up to n lines and nil on success, or a partial result
//	and a non-nil error on failure.
//
// Example:
//
//	lines, err := sysx.Head("/var/log/app.log", 10)
//	for _, l := range lines {
//	    fmt.Println(l)
//	}
func Head(path string, n int) ([]string, error) {
	if n < 0 {
		n = 0
	}
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	lines := make([]string, 0, n)
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		if len(lines) >= n {
			break
		}
		lines = append(lines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return lines, err
	}
	return lines, nil
}

// Tail reads the last n lines of the file at path and returns them as a
// slice of strings with line endings stripped. If the file has fewer than n
// lines, all lines are returned.
//
// The function is memory-efficient: it seeks towards the end of the file and
// reads only the necessary trailing portion.
//
// Parameters:
//   - `path`: the file system path to read.
//   - `n`:    the maximum number of lines to return.
//
// Returns:
//
//	([]string, error): up to n lines and nil on success, or nil and a non-nil
//	error if the file cannot be opened or read.
//
// Example:
//
//	logs, err := sysx.Tail("/var/log/app.log", 20)
//	for _, line := range logs {
//	    fmt.Println(line)
//	}
func Tail(path string, n int) ([]string, error) {
	if n <= 0 {
		return []string{}, nil
	}

	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	stat, err := f.Stat()
	if err != nil {
		return nil, err
	}

	size := stat.Size()
	var cursor int64 = size
	var lineCount int
	var buf []byte
	chunkSize := int64(4096)

	for cursor > 0 && lineCount <= n {
		if cursor < chunkSize {
			chunkSize = cursor
		}
		cursor -= chunkSize
		_, err := f.Seek(cursor, io.SeekStart)
		if err != nil {
			return nil, err
		}

		chunk := make([]byte, chunkSize)
		_, err = f.Read(chunk)
		if err != nil {
			return nil, err
		}

		buf = append(chunk, buf...)
		lineCount = bytes.Count(buf, []byte{'\n'})
	}

	allLines := splitLines(string(buf))
	if len(allLines) <= n {
		return allLines, nil
	}
	return allLines[len(allLines)-n:], nil
}
