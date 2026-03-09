package sysx

import (
	"bufio"
	"io"
	"os"
)

// CountLines counts the number of lines in the file at path.
//
// Counting is performed with a buffered scanner, making it memory-efficient
// for large files: only one line is held in memory at a time. An empty file
// returns 0 with a nil error.
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
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		count++
	}
	if err := scanner.Err(); err != nil {
		return count, err
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

// CopyFile copies the contents of the file at src to a new or truncated file
// at dst. The destination file is created with permission 0644 if it does not
// exist, or truncated if it does. The copy is performed using io.Copy with a
// buffered writer for efficiency.
//
// Parameters:
//   - `src`: the source file path to read from.
//   - `dst`: the destination file path to write to.
//
// Returns:
//
//	An error if the source cannot be opened, the destination cannot be
//	created, or the copy fails; nil on success.
//
// Example:
//
//	if err := sysx.CopyFile("/etc/hosts", "/tmp/hosts.bak"); err != nil {
//	    log.Fatal(err)
//	}
func CopyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()
	out, err := os.OpenFile(dst, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0o644)
	if err != nil {
		return err
	}
	defer out.Close()
	w := bufio.NewWriter(out)
	if _, err := io.Copy(w, in); err != nil {
		return err
	}
	return w.Flush()
}

// TruncateFile truncates or extends the file at path to the given size in
// bytes. If size is greater than the current file size, the file is extended
// with zero bytes. If the file does not exist, an error is returned.
//
// Parameters:
//   - `path`: the file system path to truncate.
//   - `size`: the desired file size in bytes; must be non-negative.
//
// Returns:
//
//	An error if the file does not exist or the truncation fails; nil on success.
//
// Example:
//
//	if err := sysx.TruncateFile("/tmp/output.bin", 0); err != nil {
//	    log.Fatal(err)
//	}
func TruncateFile(path string, size int64) error {
	return os.Truncate(path, size)
}
