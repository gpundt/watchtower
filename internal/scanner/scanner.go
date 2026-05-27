package scanner

import(
	"errors"
	"fmt"
	"net"
	"net/netip"

	ICMP "watchtower/internal/scanner/icmp"

	"github.com/rs/zerolog/log"
)

func InitializeNetworkScanner() {
	subnets, err := getHostIPSubnets()
	if err != nil {
		log.Err(err).Str("func", "getHostIPSubnets").Msg("")
		return
	}
	log.Info().Msg("Network Scanner: Initialized")

	StartScanning(subnets)
}

func StartScanning(subnets []string) {
	ICMP.RunICMPScan(subnets)
}

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
				// Parse subnet address from CIDR
				cidr := ipNet.String()
				prefix, err := netip.ParsePrefix(cidr)
				if err != nil {
					log.Err(err).Str("func", "getHostIPSubnets").Msg(fmt.Sprintf("Failed to parse prefix from '%s'", cidr))
					continue
				}
				networkAddr := string(prefix.Masked().String())
				
				// Append each network address to subnets slice
				subnets = append(subnets, networkAddr)
			}
		}
	}

	if len(subnets) == 0 {
		return nil, errors.New("No subnets found")

	} else {
		log.Debug().Str("func", "getHostIPSubnets").Msg(fmt.Sprintf("%s", subnets))
	}
	
	return subnets, nil
}