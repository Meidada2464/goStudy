/**
 * Package util
 * @Author fuqiang.li <fuqiang.li@baishan.com>
 * @Date 2024/3/7
 */

package util

import "net"

// GetNICs gets all NIC name and unicast IP address, skip the NICs that are not up and loop back.
//
// Returns a map of IP address to NIC name.
func GetNICs() (map[string]string, error) {
	// Get a list of all network interfaces.
	interfaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}

	var nics = make(map[string]string)

	for _, intf := range interfaces {
		// Check if the interface is up and skip if not.
		if intf.Flags&net.FlagUp == 0 {
			continue
		}

		// Skip loopback interfaces for unicast addresses.
		if intf.Flags&net.FlagLoopback != 0 {
			continue
		}

		addrs, err := intf.Addrs()
		if err != nil {
			continue
		}

		for _, addr := range addrs {
			// Using type assertion to convert net.Addr to *net.IPNet to access the IP address.
			ipNet, ok := addr.(*net.IPNet)
			if !ok {
				continue
			}

			// Filter out non-unicast addresses (like multicast and loopback)
			if ipNet.IP.IsLoopback() || ipNet.IP.IsMulticast() {
				continue
			}

			nics[ipNet.IP.String()] = intf.Name
		}
	}
	return nics, nil
}
