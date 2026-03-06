package netx

import (
	"math/big"
	"net"
)

// cloneIP returns a copy of an IP address, always in 4-byte (IPv4) or
// 16-byte (IPv6) form, consistent with what net.ParseCIDR produces.
func cloneIP(ip net.IP) net.IP {
	clone := make(net.IP, len(ip))
	copy(clone, ip)
	return clone
}

// ipAddOffset adds a signed offset to an IP address and returns the result.
// The operation is performed using big.Int arithmetic so that it works
// correctly for both IPv4 and IPv6 addresses without overflow.
//
// Negative offsets subtract from the address. The caller must ensure the
// result stays within the valid address range.
func ipAddOffset(ip net.IP, offset int) net.IP {
	b := new(big.Int).SetBytes(ip)
	b.Add(b, big.NewInt(int64(offset)))
	raw := b.Bytes()
	// Pad to original length so the IP type (4 vs 16 bytes) is preserved.
	result := make(net.IP, len(ip))
	copy(result[len(ip)-len(raw):], raw)
	return result
}

// ipToInt converts an IP address to a big.Int for arithmetic comparison.
func ipToInt(ip net.IP) *big.Int {
	return new(big.Int).SetBytes(ip)
}

// intToIP converts a big.Int back to a net.IP of the given byte length.
func intToIP(n *big.Int, length int) net.IP {
	raw := n.Bytes()
	ip := make(net.IP, length)
	copy(ip[length-len(raw):], raw)
	return ip
}

// computeBroadcast returns the broadcast address by ORing the network address
// with the bitwise complement of the mask.
//
// Works for both IPv4 (4-byte) and IPv6 (16-byte) addresses.
func computeBroadcast(network net.IP, mask net.IPMask) net.IP {
	broadcast := make(net.IP, len(network))
	for i := range network {
		broadcast[i] = network[i] | ^mask[i]
	}
	return broadcast
}

// computeHostRange returns the first and last usable host addresses.
//
// Edge cases:
//   - prefix == bits (e.g. /32 for IPv4, /128 for IPv6): single host; first == last == network.
//   - prefix == bits-1 (e.g. /31 for IPv4): RFC 3021 point-to-point; first == network, last == broadcast.
//   - all other prefixes: first is network+1, last is broadcast-1.
func computeHostRange(network, broadcast net.IP, prefix, bits int) (first, last net.IP) {
	switch {
	case prefix == bits:
		// /32 or /128 — single host address
		first = cloneIP(network)
		last = cloneIP(network)
	case prefix == bits-1:
		// /31 or /127 — RFC 3021 point-to-point link
		first = cloneIP(network)
		last = cloneIP(broadcast)
	default:
		first = ipAddOffset(network, 1)
		last = ipAddOffset(broadcast, -1)
	}
	return first, last
}

// computeTotalHosts returns the number of usable host addresses.
//
// Rules:
//   - prefix == bits        → 1
//   - prefix == bits-1      → 2
//   - otherwise             → 2^(bits-prefix) - 2
func computeTotalHosts(prefix, bits int) *big.Int {
	hostBits := bits - prefix
	switch {
	case hostBits == 0:
		return big.NewInt(1)
	case hostBits == 1:
		return big.NewInt(2)
	default:
		total := new(big.Int).Lsh(big.NewInt(1), uint(hostBits)) // 2^hostBits
		total.Sub(total, big.NewInt(2))                          // subtract network and broadcast
		return total
	}
}
