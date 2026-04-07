package sysx

import "os"

const (
	// os.O_CREATE | os.O_WRONLY | os.O_APPEND
	// Flag for opening a file for writing with append mode and creating it if it doesn't exist.
	CWA FileOpenFlags = FileOpenFlags(os.O_CREATE | os.O_WRONLY | os.O_APPEND)

	// os.O_CREATE | os.O_WRONLY | os.O_TRUNC
	// Flag for opening a file for writing with truncate mode and creating it if it doesn't exist.
	CWT FileOpenFlags = FileOpenFlags(os.O_CREATE | os.O_WRONLY | os.O_TRUNC)

	// os.O_CREATE | os.O_WRONLY
	// Flag for opening a file for writing and creating it if it doesn't exist.
	CW FileOpenFlags = FileOpenFlags(os.O_CREATE | os.O_WRONLY)

	// os.O_APPEND | os.O_WRONLY
	// Flag for opening a file for writing with append mode.
	AW FileOpenFlags = FileOpenFlags(os.O_APPEND | os.O_WRONLY)

	// os.O_TRUNC | os.O_WRONLY
	// Flag for opening a file for writing with truncate mode.
	TW FileOpenFlags = FileOpenFlags(os.O_TRUNC | os.O_WRONLY)

	// os.O_RDONLY
	// Flag for opening a file for reading only.
	RO FileOpenFlags = FileOpenFlags(os.O_RDONLY)

	// os.O_RDWR
	// Flag for opening a file for reading and writing.
	RW FileOpenFlags = FileOpenFlags(os.O_RDWR)

	// os.O_RDWR | os.O_CREATE
	// Flag for opening a file for reading and writing and creating it if it doesn't exist.
	CRW FileOpenFlags = FileOpenFlags(os.O_RDWR | os.O_CREATE)

	// os.O_RDWR | os.O_CREATE | os.O_TRUNC
	// Flag for opening a file for reading and writing with truncate mode and creating it if it doesn't exist.
	CRWT FileOpenFlags = FileOpenFlags(os.O_RDWR | os.O_CREATE | os.O_TRUNC)

	// os.O_RDWR | os.O_APPEND
	// Flag for opening a file for reading and writing with append mode.
	ARW FileOpenFlags = FileOpenFlags(os.O_RDWR | os.O_APPEND)

	// os.O_RDWR | os.O_APPEND | os.O_TRUNC
	// Flag for opening a file for reading and writing with append and truncate mode and creating it if it doesn't exist.
	TARW FileOpenFlags = FileOpenFlags(os.O_RDWR | os.O_APPEND | os.O_TRUNC)

	// os.O_RDWR | os.O_TRUNC
	// Flag for opening a file for reading and writing with truncate mode.
	TRW FileOpenFlags = FileOpenFlags(os.O_RDWR | os.O_TRUNC)

	// os.O_RDWR | os.O_CREATE | os.O_APPEND
	// Flag for opening a file for reading and writing with append mode and creating it if it doesn't exist.
	CRWA FileOpenFlags = FileOpenFlags(os.O_RDWR | os.O_CREATE | os.O_APPEND)
)
