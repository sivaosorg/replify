package randn

// XID represents a unique request id, cloned and improved from rs/xid.
// It is a 12-byte identifier composed of a 4-byte timestamp, 3-byte machine ID,
// 2-byte process ID, and a 3-byte counter.
type XID [12]byte
