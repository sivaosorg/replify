package netx

import "net"

// buildSubnet computes all addressing attributes from an *net.IPNet and
// returns a fully populated Subnet value.
//
// buildSubnet is the single place where all address arithmetic happens,
// making it the foundation for both ParseCIDR and the FLSM/VLSM allocators.
func buildSubnet(ipNet *net.IPNet) Subnet {
	prefix, bits := ipNet.Mask.Size()
	network := cloneIP(ipNet.IP)
	broadcast := computeBroadcast(network, ipNet.Mask)
	first, last := computeHostRange(network, broadcast, prefix, bits)
	total := computeTotalHosts(prefix, bits)

	return Subnet{
		ipNet:            ipNet,
		networkAddress:   network,
		broadcastAddress: broadcast,
		firstHost:        first,
		lastHost:         last,
		totalHosts:       total,
		prefix:           prefix,
	}
}
