package scanner

import (
	"errors"
	"fmt"
	"net"

	//ARP "watchtower/internal/scanner/arp"
	ICMP "watchtower/internal/scanner/icmp"
	Port "watchtower/internal/scanner/port"

	"github.com/rs/zerolog/log"
)

// Function to run network scan on all subnets the server is connected to
func StartNetworkScanner() {
	subnets, err := getHostIPSubnets()
	if err != nil {
		log.Err(err).Str("func", "getHostIPSubnets").Msg("")
		return
	}
	log.Debug().Msg("Network Scanner: Initialized")

	ICMP.RunICMPScan(subnets)
	//ARP.RunARPScan(subnets)
	Port.RunPortScan(subnets)
}

// Function to get and return all subnets the server is connected to 
func getHostIPSubnets() (map[string]string, error) {
	subnets := map[string]string{}

	interfaces, err := net.Interfaces()
	if err != nil {
		return nil, errors.New("Failed to get net.Interfaces()")
	}

	// For each interface in net.Interfaces
	for _, iface := range interfaces {
		addrs, err := iface.Addrs()
		if err != nil {
			return nil, errors.New(fmt.Sprintf(
				"Failed to fetch addresses for %s: %v",
				iface.Name,
				err,
			))
		}

		// For each IP on each interface
		for _, addr := range addrs {
			if ipNet, ok := addr.(*net.IPNet); ok && !ipNet.IP.IsLoopback() {
				if ipNet.IP.To4() != nil {
					// Get each address; /24 subnet
					ipv4 := ipNet.IP.To4()
					normalized := fmt.Sprintf(
						"%d.%d.%d.0/24",
						ipv4[0],
						ipv4[1],
						ipv4[2],
					)

					// Map each interface to its subnet
					subnets[iface.Name] = normalized
				}
			}
		}
	}

	if len(subnets) == 0 {
		return nil, errors.New("No subnets found")

	} else {
		log.Debug().
			Str("func", "getHostIPSubnets").
			Msg(fmt.Sprintf("%s", subnets))
	}

	return subnets, nil
}
