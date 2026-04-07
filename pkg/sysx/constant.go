package sysx

import "os"

// FileOpenFlags is a convenience type for file open flags.
const (
	// os.O_CREATE | os.O_WRONLY | os.O_APPEND
	CWA = os.O_CREATE | os.O_WRONLY | os.O_APPEND

	// os.O_CREATE | os.O_WRONLY | os.O_TRUNC
	CWT = os.O_CREATE | os.O_WRONLY | os.O_TRUNC

	// os.O_CREATE | os.O_WRONLY
	CW = os.O_CREATE | os.O_WRONLY

	// os.O_APPEND | os.O_WRONLY
	AW = os.O_APPEND | os.O_WRONLY

	// os.O_TRUNC | os.O_WRONLY
	TW = os.O_TRUNC | os.O_WRONLY

	// os.O_RDONLY
	RO = os.O_RDONLY

	// os.O_RDWR
	RW = os.O_RDWR

	// os.O_RDWR | os.O_CREATE
	CRW = os.O_RDWR | os.O_CREATE

	// os.O_RDWR | os.O_CREATE | os.O_TRUNC
	CRWT = os.O_RDWR | os.O_CREATE | os.O_TRUNC

	// os.O_RDWR | os.O_APPEND
	ARW = os.O_RDWR | os.O_APPEND

	// os.O_RDWR | os.O_APPEND | os.O_CREATE
	CARW = os.O_RDWR | os.O_APPEND | os.O_CREATE

	// os.O_RDWR | os.O_APPEND | os.O_TRUNC
	TARW = os.O_RDWR | os.O_APPEND | os.O_TRUNC

	// os.O_RDWR | os.O_TRUNC
	TRW = os.O_RDWR | os.O_TRUNC

	// os.O_RDWR | os.O_CREATE | os.O_APPEND
	CRWA = os.O_RDWR | os.O_CREATE | os.O_APPEND
)
