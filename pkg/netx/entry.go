package netx

import (
	"errors"
	"fmt"
	"math/big"
	"net"
)

// Contains reports whether the given IP address belongs to the network.
//
// This is a thin, descriptive wrapper around (*net.IPNet).Contains.
//
// Parameters:
//   - `network`: the network to test against.
//   - `ip`:      the IP address to look up.
//
// Returns:
//
//	A boolean value:
//	 - true  when ip falls within the network's address range;
//	 - false otherwise or when either argument is nil.
//
// Example:
//
//	_, n, _ := net.ParseCIDR("10.0.0.0/8")
//	netx.Contains(n, net.ParseIP("10.1.2.3"))  // true
//	netx.Contains(n, net.ParseIP("192.168.1.1")) // false
func Contains(network *net.IPNet, ip net.IP) bool {
	if network == nil || ip == nil {
		return false
	}
	return network.Contains(ip)
}

// Overlaps reports whether two network blocks share any address in common.
//
// Two networks overlap when one contains the network address of the other,
// or when one is entirely contained within the other.
//
// Parameters:
//   - `netA`: the first network.
//   - `netB`: the second network.
//
// Returns:
//
//	A boolean value:
//	 - true  when the two networks share at least one address;
//	 - false when they are disjoint or either argument is nil.
//
// Example:
//
//	_, a, _ := net.ParseCIDR("10.0.0.0/24")
//	_, b, _ := net.ParseCIDR("10.0.0.128/25")
//	netx.Overlaps(a, b) // true
//
//	_, c, _ := net.ParseCIDR("192.168.0.0/24")
//	netx.Overlaps(a, c) // false
func Overlaps(netA, netB *net.IPNet) bool {
	if netA == nil || netB == nil {
		return false
	}
	return netA.Contains(netB.IP) || netB.Contains(netA.IP)
}

// NetworkSize returns the total number of addresses in the network block
// (including network and broadcast addresses). For a /24 this is 256.
//
// A *big.Int is used so that IPv6 networks, which can contain up to 2^128
// addresses, are handled correctly.
//
// Parameters:
//   - `ipnet`: the network whose size to compute.
//
// Returns:
//
//	A *big.Int representing the total address count; 0 when ipnet is nil.
//
// Example:
//
//	_, n, _ := net.ParseCIDR("192.168.1.0/24")
//	netx.NetworkSize(n) // 256
func NetworkSize(ipnet *net.IPNet) *big.Int {
	if ipnet == nil {
		return big.NewInt(0)
	}
	prefix, bits := ipnet.Mask.Size()
	return new(big.Int).Lsh(big.NewInt(1), uint(bits-prefix))
}

// HostCount returns the number of usable host addresses for a subnet with the
// given prefix length in a standard IPv4 (/8–/32) or IPv6 (/0–/128) network.
//
// Special cases:
//   - prefix == bits (e.g. /32 for IPv4): returns 1.
//   - prefix == bits-1 (e.g. /31 for IPv4): returns 2.
//   - all other prefixes: returns 2^(bits-prefix) - 2.
//
// bits must be either 32 (IPv4) or 128 (IPv6). If bits is neither, the
// function treats the prefix as IPv4.
//
// Parameters:
//   - `prefix`: the subnet prefix length.
//   - `bits`:   the total number of bits in the address family (32 or 128).
//
// Returns:
//
//	A *big.Int representing the usable host count.
//
// Example:
//
//	netx.HostCount(24, 32) // 254
//	netx.HostCount(31, 32) // 2
//	netx.HostCount(32, 32) // 1
func HostCount(prefix, bits int) *big.Int {
	return computeTotalHosts(prefix, bits)
}

// PrefixForHosts returns the smallest prefix length that provides at least
// the requested number of usable host addresses.
//
// The function searches from the most specific prefix toward the least
// specific, returning the first prefix for which HostCount(prefix, bits) ≥
// hosts.
//
// Parameters:
//   - `hosts`: the minimum number of usable host addresses required (≥ 1).
//   - `bits`:  the address family bit width (32 for IPv4, 128 for IPv6).
//
// Returns:
//
//	An int containing the prefix length, or -1 when no valid prefix can
//	satisfy the requirement (e.g. more than 2^30 hosts requested for IPv4).
//
// Example:
//
//	netx.PrefixForHosts(100, 32)  // 25 (provides 126 usable hosts)
//	netx.PrefixForHosts(254, 32)  // 24 (provides 254 usable hosts)
//	netx.PrefixForHosts(255, 32)  // 23 (provides 510 usable hosts)
func PrefixForHosts(hosts, bits int) int {
	required := big.NewInt(int64(hosts))
	// Walk from most-specific to least-specific prefix.
	for prefix := bits; prefix >= 0; prefix-- {
		if computeTotalHosts(prefix, bits).Cmp(required) >= 0 {
			return prefix
		}
	}
	return -1
}

// NextSubnet returns the next contiguous subnet of the given prefix size that
// immediately follows ipnet.
//
// For example, the next /26 after 10.0.0.0/26 is 10.0.0.64/26. The function
// does not check whether the returned subnet is within any enclosing block.
//
// Parameters:
//   - `ipnet`:     the current subnet.
//   - `newPrefix`: the prefix length for the next subnet.
//
// Returns:
//
//	(*net.IPNet, error): the next subnet, or nil and a non-nil error when
//	ipnet is nil or newPrefix is invalid.
//
// Example:
//
//	base := netx.MustParseCIDR("10.0.0.0/26").IPNet()
//	next, err := netx.NextSubnet(base, 26)
//	fmt.Println(next) // "10.0.0.64/26"
func NextSubnet(ipnet *net.IPNet, newPrefix int) (*net.IPNet, error) {
	if ipnet == nil {
		return nil, errors.New("netx: NextSubnet: ipnet must not be nil")
	}
	_, bits := ipnet.Mask.Size()
	if newPrefix < 0 || newPrefix > bits {
		return nil, fmt.Errorf(
			"netx: NextSubnet: newPrefix (%d) out of range [0, %d]",
			newPrefix, bits,
		)
	}

	// Compute the broadcast of the current subnet, then add 1.
	broadcast := computeBroadcast(ipnet.IP, ipnet.Mask)
	nextIP := ipAddOffset(broadcast, 1)

	// Apply the requested new mask.
	newMask := net.CIDRMask(newPrefix, bits)
	// Mask to find the network address.
	nextNet := make(net.IP, len(nextIP))
	for i := range nextIP {
		nextNet[i] = nextIP[i] & newMask[i]
	}

	return &net.IPNet{IP: nextNet, Mask: newMask}, nil
}

// SubnetsToStrings converts a slice of *net.IPNet to their CIDR string
// representations.
//
// Parameters:
//   - `nets`: the slice of networks to convert.
//
// Returns:
//
//	A slice of CIDR strings in the same order as the input.
//
// Example:
//
//	base := netx.MustParseCIDR("10.0.0.0/24").IPNet()
//	subs, _ := netx.Split(base, 26)
//	fmt.Println(netx.SubnetsToStrings(subs))
//	// ["10.0.0.0/26" "10.0.0.64/26" "10.0.0.128/26" "10.0.0.192/26"]
func SubnetsToStrings(nets []*net.IPNet) []string {
	result := make([]string, len(nets))
	for i, n := range nets {
		if n != nil {
			result[i] = n.String()
		}
	}
	return result
}

// AllocatedSubnetsToStrings converts a slice of Subnet to their CIDR string
// representations.
//
// Parameters:
//   - `subnets`: the slice of Subnet values to convert.
//
// Returns:
//
//	A slice of CIDR strings in the same order as the input.
//
// Example:
//
//	base := netx.MustParseCIDR("10.0.0.0/24").IPNet()
//	subs, _ := netx.DivideByHosts(base, []int{100, 50, 10})
//	fmt.Println(netx.AllocatedSubnetsToStrings(subs))
func AllocatedSubnetsToStrings(subnets []Subnet) []string {
	result := make([]string, len(subnets))
	for i, s := range subnets {
		result[i] = s.String()
	}
	return result
}
