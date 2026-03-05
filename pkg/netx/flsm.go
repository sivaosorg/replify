package netx

import (
	"errors"
	"fmt"
	"math/big"
	"net"
)

// ///////////////////////////
// Section: Fixed-Length Subnet Masking (FLSM)
// ///////////////////////////

// Split divides a network block into equal-sized subnets, each with the given
// prefix length.
//
// The function performs Fixed-Length Subnet Masking (FLSM): all resulting
// subnets are the same size, and together they exactly cover the original
// network without gaps or overlaps.
//
// Parameters:
//   - `network`:   the base network to split.
//   - `newPrefix`: the prefix length for each resulting subnet; must be
//     strictly greater than network's prefix.
//
// Returns:
//
//	([]*net.IPNet, error): an ordered slice of subnets covering the original
//	network, or nil and a non-nil error when the split is invalid.
//
// Errors:
//   - When newPrefix is not larger than the base prefix.
//   - When newPrefix exceeds the maximum for the address family (32 for IPv4,
//     128 for IPv6).
//
// Example:
//
//	base := netx.MustParseCIDR("10.0.0.0/24").IPNet()
//	subnets, err := netx.Split(base, 26)
//	// subnets: [10.0.0.0/26, 10.0.0.64/26, 10.0.0.128/26, 10.0.0.192/26]
func Split(network *net.IPNet, newPrefix int) ([]*net.IPNet, error) {
	if network == nil {
		return nil, errors.New("netx: Split: network must not be nil")
	}
	basePrefix, bits := network.Mask.Size()
	if newPrefix <= basePrefix {
		return nil, fmt.Errorf(
			"netx: Split: newPrefix (%d) must be greater than base prefix (%d)",
			newPrefix, basePrefix,
		)
	}
	if newPrefix > bits {
		return nil, fmt.Errorf(
			"netx: Split: newPrefix (%d) exceeds maximum prefix (%d) for this address family",
			newPrefix, bits,
		)
	}

	// Number of subnets = 2^(newPrefix - basePrefix)
	extraBits := uint(newPrefix - basePrefix)
	count := new(big.Int).Lsh(big.NewInt(1), extraBits)

	// Size of each subnet in addresses = 2^(bits - newPrefix)
	subnetSize := new(big.Int).Lsh(big.NewInt(1), uint(bits-newPrefix))

	newMask := net.CIDRMask(newPrefix, bits)
	ipLen := len(network.IP)
	base := ipToInt(network.IP)

	var result []*net.IPNet
	current := new(big.Int).Set(base)
	for i := new(big.Int); i.Cmp(count) < 0; i.Add(i, big.NewInt(1)) {
		ip := intToIP(current, ipLen)
		result = append(result, &net.IPNet{
			IP:   ip,
			Mask: newMask,
		})
		current.Add(current, subnetSize)
	}
	return result, nil
}

// SplitIntoN divides a network block into exactly n equal-sized subnets.
//
// SplitIntoN is a convenience wrapper around Split that automatically
// calculates the required prefix length. n must be a power of two.
//
// Parameters:
//   - `network`: the base network to split.
//   - `n`:       the number of subnets to produce; must be a power of 2 and
//     at least 2.
//
// Returns:
//
//	([]*net.IPNet, error): n equal-sized subnets, or nil and an error.
//
// Example:
//
//	base := netx.MustParseCIDR("10.0.0.0/24").IPNet()
//	subnets, err := netx.SplitIntoN(base, 4)
//	// subnets: [10.0.0.0/26, 10.0.0.64/26, 10.0.0.128/26, 10.0.0.192/26]
func SplitIntoN(network *net.IPNet, n int) ([]*net.IPNet, error) {
	if network == nil {
		return nil, errors.New("netx: SplitIntoN: network must not be nil")
	}
	if n < 2 {
		return nil, fmt.Errorf("netx: SplitIntoN: n (%d) must be at least 2", n)
	}
	if n&(n-1) != 0 {
		return nil, fmt.Errorf("netx: SplitIntoN: n (%d) must be a power of 2", n)
	}
	basePrefix, _ := network.Mask.Size()
	bitsNeeded := 0
	v := n
	for v > 1 {
		v >>= 1
		bitsNeeded++
	}
	return Split(network, basePrefix+bitsNeeded)
}
