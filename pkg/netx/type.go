package netx

import (
	"math/big"
	"net"
)

// Subnet represents a fully computed IP network block, including all derived
// addressing attributes. It is the central data structure of the netx package.
//
// A Subnet is created by ParseCIDR or returned from FLSM/VLSM allocation
// functions. All fields are computed automatically; use the accessor methods
// to read them.
//
// Subnet is safe for concurrent reads after creation. It must not be modified
// after construction.
type Subnet struct {
	// ipNet is the parsed network from Go's net package.
	ipNet *net.IPNet

	// networkAddress is the lowest address in the block (all host bits zero).
	networkAddress net.IP

	// broadcastAddress is the highest address in the block (all host bits one).
	// For IPv6 and /31 or /32, this follows the same bitwise rule.
	broadcastAddress net.IP

	// firstHost is the first usable host address.
	// For /31 networks, this equals networkAddress.
	// For /32 networks, this equals the single host address.
	firstHost net.IP

	// lastHost is the last usable host address.
	// For /31 networks, this equals broadcastAddress.
	// For /32 networks, this equals the single host address.
	lastHost net.IP

	// totalHosts is the number of usable host addresses in the block.
	// /31 returns 2 (point-to-point), /32 returns 1.
	totalHosts *big.Int

	// prefix is the subnet prefix length (e.g. 24 for a /24).
	prefix int
}
