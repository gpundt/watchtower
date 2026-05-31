package scanner

import (
	"errors"
	"fmt"
	"net"
	"slices"

	ARP "watchtower/internal/scanner/arp"
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
	ARP.RunARPScan()
	Port.RunPortScan(subnets)
}

// Function to get and return all subnets the server is connected to 
func getHostIPSubnets() ([]string, error) {
	subnets := []string{}

	// Isolate addresses on this host
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return nil, errors.New("Failed to get net.InterfaceAddrs")
	}

	// Iterate through the slice of addresses
	for _, address := range addrs {
		if ipNet, ok := address.(*net.IPNet); ok && !ipNet.IP.IsLoopback() {
			if ipNet.IP.To4() != nil {
				// Get each address; /24 subnet
				ipv4 := ipNet.IP.To4()
				normalized := fmt.Sprintf(
					"%d.%d.%d.0/24",
					ipv4[0],
					ipv4[1],
					ipv4[2],
				)

				// Append each network address to subnets slice
				if !slices.Contains(subnets, normalized) {
					subnets = append(subnets, normalized)
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
