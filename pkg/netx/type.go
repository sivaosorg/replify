package netx

import (
	"math/big"
	"net"
)

// ///////////////////////////
// Section: Subnet type
// ///////////////////////////

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

// IPNet returns the underlying *net.IPNet for this subnet.
//
// Returns:
//
//	A pointer to the net.IPNet representing this network block.
func (s Subnet) IPNet() *net.IPNet { return s.ipNet }

// NetworkAddress returns the network (base) address of the subnet.
//
// Returns:
//
//	A net.IP containing the lowest address of the block.
//
// Example:
//
//	sub, _ := netx.ParseCIDR("10.0.0.0/24")
//	fmt.Println(sub.NetworkAddress()) // "10.0.0.0"
func (s Subnet) NetworkAddress() net.IP { return s.networkAddress }

// BroadcastAddress returns the broadcast address of the subnet.
//
// For IPv6 subnets and /31 or /32 blocks the value is still the
// bitwise all-ones host address, even though broadcast semantics differ.
//
// Returns:
//
//	A net.IP containing the highest address of the block.
//
// Example:
//
//	sub, _ := netx.ParseCIDR("10.0.0.0/24")
//	fmt.Println(sub.BroadcastAddress()) // "10.0.0.255"
func (s Subnet) BroadcastAddress() net.IP { return s.broadcastAddress }

// FirstHost returns the first usable host address.
//
// For /31 (RFC 3021) and /32 blocks the concept of "first host" maps to
// networkAddress and the single host address respectively.
//
// Returns:
//
//	A net.IP containing the first usable address.
//
// Example:
//
//	sub, _ := netx.ParseCIDR("10.0.0.0/24")
//	fmt.Println(sub.FirstHost()) // "10.0.0.1"
func (s Subnet) FirstHost() net.IP { return s.firstHost }

// LastHost returns the last usable host address.
//
// For /31 and /32 blocks see the note on FirstHost.
//
// Returns:
//
//	A net.IP containing the last usable address.
//
// Example:
//
//	sub, _ := netx.ParseCIDR("10.0.0.0/24")
//	fmt.Println(sub.LastHost()) // "10.0.0.254"
func (s Subnet) LastHost() net.IP { return s.lastHost }

// TotalHosts returns the number of usable host addresses in the subnet.
//
// A *big.Int is returned to handle IPv6 subnets whose host counts exceed
// int64 range. For typical IPv4 subnets the value fits in int64.
//
// Special cases:
//   - /31 returns 2 (RFC 3021 point-to-point link)
//   - /32 returns 1
//
// Returns:
//
//	A *big.Int representing the usable host count.
//
// Example:
//
//	sub, _ := netx.ParseCIDR("10.0.0.0/24")
//	fmt.Println(sub.TotalHosts()) // 254
func (s Subnet) TotalHosts() *big.Int { return new(big.Int).Set(s.totalHosts) }

// Prefix returns the prefix length of the subnet (e.g. 24 for a /24).
//
// Returns:
//
//	An int containing the prefix length.
//
// Example:
//
//	sub, _ := netx.ParseCIDR("192.168.1.0/24")
//	fmt.Println(sub.Prefix()) // 24
func (s Subnet) Prefix() int { return s.prefix }

// String returns the CIDR notation of the subnet (e.g. "10.0.0.0/24").
//
// Returns:
//
//	A string in CIDR notation.
func (s Subnet) String() string { return s.ipNet.String() }
