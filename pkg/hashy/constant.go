package hashy

const (
	// visitFlagSet is a flag that indicates that the slice should be treated as a set (order-independent)
	visitFlagSet visitFlag = 1 << iota
)

const (
	// H_CRC32 is a CRC32 hash algorithm.
	H_CRC32 HashAlgorithm = "crc32"

	// H_CRC64 is a CRC64 hash algorithm.
	H_CRC64 HashAlgorithm = "crc64"

	// H_MD5 is a MD5 hash algorithm.
	H_MD5 HashAlgorithm = "md5"

	// H_SHA1 is a SHA1 hash algorithm.
	H_SHA1 HashAlgorithm = "sha1"

	// H_SHA224 is a SHA224 hash algorithm.
	H_SHA224 HashAlgorithm = "sha224"

	// H_SHA256 is a SHA256 hash algorithm.
	H_SHA256 HashAlgorithm = "sha256"

	// H_SHA384 is a SHA384 hash algorithm.
	H_SHA384 HashAlgorithm = "sha384"

	// H_SHA512 is a SHA512 hash algorithm.
	H_SHA512 HashAlgorithm = "sha512"
)
