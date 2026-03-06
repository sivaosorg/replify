package netx

import (
	"fmt"
	"net"
)

// ParseCIDR parses a CIDR notation string and returns a fully populated
// Subnet with all addressing attributes calculated.
//
// The CIDR string must be in standard notation, for example "192.168.1.0/24"
// or "2001:db8::/32". Both IPv4 and IPv6 are supported.
//
// Unlike net.ParseCIDR, which silently masks the host bits, ParseCIDR always
// uses the network address derived from the mask — making it safe to pass
// host addresses such as "192.168.1.5/24" and receive the correct network.
//
// Parameters:
//   - `cidr`: a CIDR notation string (e.g. "10.0.0.0/8").
//
// Returns:
//
//	(Subnet, error): a fully computed Subnet on success, or a zero Subnet
//	and a non-nil error when the input is malformed.
//
// Example:
//
//	sub, err := netx.ParseCIDR("192.168.1.0/24")
//	if err != nil {
//	    log.Fatal(err)
//	}
//	fmt.Println(sub.NetworkAddress())   // 192.168.1.0
//	fmt.Println(sub.BroadcastAddress()) // 192.168.1.255
//	fmt.Println(sub.FirstHost())        // 192.168.1.1
//	fmt.Println(sub.LastHost())         // 192.168.1.254
//	fmt.Println(sub.TotalHosts())       // 254
func ParseCIDR(cidr string) (Subnet, error) {
	_, ipNet, err := net.ParseCIDR(cidr)
	if err != nil {
		return Subnet{}, fmt.Errorf("netx: ParseCIDR(%q): %w", cidr, err)
	}
	return buildSubnet(ipNet), nil
}

// MustParseCIDR is like ParseCIDR but panics when the CIDR string is invalid.
//
// It is intended for use in tests and program initialisation where an invalid
// CIDR is a programming error rather than a runtime condition.
//
// Parameters:
//   - `cidr`: a CIDR notation string (e.g. "10.0.0.0/8").
//
// Returns:
//
//	A fully computed Subnet.
//
// Example:
//
//	sub := netx.MustParseCIDR("10.0.0.0/8")
//	fmt.Println(sub.NetworkAddress()) // 10.0.0.0
func MustParseCIDR(cidr string) Subnet {
	s, err := ParseCIDR(cidr)
	if err != nil {
		panic("netx: MustParseCIDR: " + err.Error())
	}
	return s
}
