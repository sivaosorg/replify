package randn

import (
	"errors"
	"os"
)

var (
	// ErrInvalidID is returned when trying to unmarshal an invalid ID
	ErrInvalidID = errors.New("randn: invalid XID")
)

const (
	// encodedLen is the length of the encoded XID.
	encodedLen = 20

	// encoding stores a custom version of the base32 encoding with lower case letters.
	encoding = "0123456789abcdefghijklmnopqrstuv"
)

var (
	// idCounter is a random number used to generate unique IDs.
	idCounter = RandUint32()

	// machineID is the hardware address of the machine.
	machineID = readMachineID()

	// pid is the process ID of the current process.
	pid = os.Getpid()

	// zeroXID is the zero value for XID.
	zeroXID XID

	// decodeBytes is a lookup table for decoding base32 encoded strings.
	decodeBytes [256]byte
)
