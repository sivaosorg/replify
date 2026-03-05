// Package netx provides a production-ready IPv4 and IPv6 network subnetting
// toolkit built exclusively on the Go standard library.
//
// netx is designed for backend infrastructure services, DevOps tooling,
// Kubernetes networking automation, and cloud infrastructure management where
// programmatic subnet calculation, allocation, and validation are required.
//
// # Data Structures
//
// The central type is Subnet, which captures all addressing attributes of a
// network block:
//
//   - NetworkAddress   — the lowest address (all host bits zero)
//   - BroadcastAddress — the highest address (all host bits one)
//   - FirstHost        — first usable host address
//   - LastHost         — last usable host address
//   - TotalHosts       — number of usable addresses (*big.Int, handles IPv6)
//   - Prefix           — subnet prefix length (e.g. 24 for a /24)
//
// All fields are unexported; read them using accessor methods.
//
// # CIDR Parsing
//
//	sub, err := netx.ParseCIDR("192.168.1.0/24")
//	// sub.NetworkAddress()   → 192.168.1.0
//	// sub.BroadcastAddress() → 192.168.1.255
//	// sub.FirstHost()        → 192.168.1.1
//	// sub.LastHost()         → 192.168.1.254
//	// sub.TotalHosts()       → 254
//	// sub.Prefix()           → 24
//
// # Fixed-Length Subnet Masking (FLSM)
//
// Split divides a network into equal-sized subnets:
//
//	base := netx.MustParseCIDR("10.0.0.0/24").IPNet()
//	subnets, err := netx.Split(base, 26)
//	// → [10.0.0.0/26, 10.0.0.64/26, 10.0.0.128/26, 10.0.0.192/26]
//
//	// Or split into exactly N equal parts (N must be a power of 2):
//	subnets, err = netx.SplitIntoN(base, 4)
//
// # Variable-Length Subnet Masking (VLSM)
//
// DivideByHosts allocates subnets sized to individual host requirements:
//
//	base := netx.MustParseCIDR("10.0.0.0/24").IPNet()
//	subnets, err := netx.DivideByHosts(base, []int{100, 50, 10})
//	// subnets[0]: 10.0.0.0/25   (126 hosts — satisfies 100)
//	// subnets[1]: 10.0.0.128/26 ( 62 hosts — satisfies  50)
//	// subnets[2]: 10.0.0.192/28 ( 14 hosts — satisfies  10)
//
// Requirements are automatically sorted largest-first to minimise waste.
//
// # Utility Functions
//
//	netx.Contains(network, ip)    // true when ip is within network
//	netx.Overlaps(netA, netB)     // true when two networks share addresses
//	netx.NetworkSize(ipnet)       // total addresses (*big.Int, incl. network+broadcast)
//	netx.HostCount(prefix, bits)  // usable hosts for a given prefix
//	netx.PrefixForHosts(n, bits)  // smallest prefix providing ≥ n usable hosts
//	netx.NextSubnet(ipnet, pfx)   // next contiguous subnet of given prefix
//
// # Edge Case Handling
//
// The package correctly handles:
//   - /31 (RFC 3021 point-to-point): TotalHosts = 2, FirstHost = NetworkAddress
//   - /32 single-host: TotalHosts = 1, FirstHost = LastHost = NetworkAddress
//   - IPv6 subnets of any prefix length
//   - VLSM allocation failure when address space is exhausted
//
// # Cross-Platform Compatibility
//
// netx relies only on the Go standard library (net, math/big, sort) and
// produces identical results on Linux, macOS, and Windows.
package netx
