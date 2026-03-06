package netx

import (
	"errors"
	"fmt"
	"math/big"
	"net"
	"sort"
)

// DivideByHosts allocates variable-length subnets from a base network to
// satisfy a list of host-count requirements, using Variable Length Subnet
// Masking (VLSM).
//
// The algorithm:
//  1. Sorts host requirements in descending order so that the largest subnets
//     are allocated first (minimising wasted address space).
//  2. For each requirement, determines the smallest prefix that provides at
//     least the requested number of usable hosts.
//  3. Allocates subnets sequentially from the base network address with no
//     gaps.
//  4. Returns an error if the base network does not have enough address space
//     to satisfy all requirements.
//
// Parameters:
//   - `base`:             the network block from which subnets are allocated.
//   - `hostRequirements`: the number of usable hosts required for each subnet.
//     Values must be ≥ 1.
//
// Returns:
//
//	([]Subnet, error): allocated subnets in the order corresponding to the
//	sorted requirements, or nil and a non-nil error when allocation fails.
//
// Example:
//
//	base := netx.MustParseCIDR("10.0.0.0/24").IPNet()
//	subnets, err := netx.DivideByHosts(base, []int{100, 50, 10})
//	// subnets[0]: 10.0.0.0/25   (126 usable hosts)
//	// subnets[1]: 10.0.0.128/26 ( 62 usable hosts)
//	// subnets[2]: 10.0.0.192/28 ( 14 usable hosts)
func DivideByHosts(base *net.IPNet, hostRequirements []int) ([]Subnet, error) {
	if base == nil {
		return nil, errors.New("netx: DivideByHosts: base network must not be nil")
	}
	if len(hostRequirements) == 0 {
		return nil, errors.New("netx: DivideByHosts: hostRequirements must not be empty")
	}

	_, bits := base.Mask.Size()

	// Validate all requirements before allocating.
	for i, h := range hostRequirements {
		if h < 1 {
			return nil, fmt.Errorf(
				"netx: DivideByHosts: hostRequirements[%d] = %d; must be ≥ 1", i, h)
		}
	}

	// Sort descending so largest blocks are allocated first.
	sorted := make([]int, len(hostRequirements))
	copy(sorted, hostRequirements)
	sort.Sort(sort.Reverse(sort.IntSlice(sorted)))

	// Walk through each requirement, compute the required prefix, and carve
	// out the next available subnet from the current position in the base.
	ipLen := len(base.IP)
	current := ipToInt(base.IP) // next available network address
	baseEnd := ipToInt(computeBroadcast(base.IP, base.Mask))

	result := make([]Subnet, 0, len(sorted))

	for i, hosts := range sorted {
		prefix := PrefixForHosts(hosts, bits)
		if prefix < 0 {
			return nil, fmt.Errorf(
				"netx: DivideByHosts: cannot satisfy requirement of %d hosts (no valid prefix)", hosts)
		}

		// Align current to the subnet boundary for this prefix.
		subnetSize := new(big.Int).Lsh(big.NewInt(1), uint(bits-prefix))
		// Round up current to the next multiple of subnetSize.
		rem := new(big.Int).Mod(current, subnetSize)
		if rem.Sign() != 0 {
			current.Add(current, new(big.Int).Sub(subnetSize, rem))
		}

		// Verify the subnet fits within the base network.
		subnetEnd := new(big.Int).Add(current, subnetSize)
		subnetEnd.Sub(subnetEnd, big.NewInt(1))
		if subnetEnd.Cmp(baseEnd) > 0 {
			return nil, fmt.Errorf(
				"netx: DivideByHosts: insufficient address space for requirement %d (%d hosts)",
				i, hosts)
		}

		ipNet := &net.IPNet{
			IP:   intToIP(current, ipLen),
			Mask: net.CIDRMask(prefix, bits),
		}
		result = append(result, buildSubnet(ipNet))

		// Advance past this subnet.
		current.Add(current, subnetSize)
	}

	return result, nil
}
