package port

import (
	"fmt"
	"net"
	"net/netip"
	"slices"
	"sync"
	"time"

	Config "watchtower/internal/config"
	Database "watchtower/internal/database"

	"github.com/rs/zerolog/log"
)

type PortScanResults struct {
	ActiveHosts   map[string][]int
	InactiveHosts []string
}

// Function to check ports 0-1024 on all hosts in the /24 subnet
func RunPortScan(subnets map[string]string) {
	log.Info().Msg("Beginning Port Scan")

	scanResults := PortScanResults{
		ActiveHosts:   map[string][]int{},
		InactiveHosts: []string{},
	}
	guard := make(
		chan struct{},
		Config.ServerConfig.Scanner.MaxConcurrentScans,
	)

	var wg sync.WaitGroup

	// Iterate over each subnet in subnets
	for _, subnet := range subnets {
		prefix, err := netip.ParsePrefix(subnet)
		if err != nil {
			log.Err(err).Str("func", "RunPortScan").Msg("")
		}

		// for i := networkStart; i <= networkEnd; i++ {
		for addr := prefix.Addr(); prefix.Contains(addr); addr = addr.Next() {
			wg.Add(1)
			guard <- struct{}{} // Block if max number of goroutines alrady running

			go func(ip string) {
				defer wg.Done()
				defer func() { <-guard }() // Release the slot in the semaphore

				hostActive := false
				activePorts := []int{}
				for port := range 1024 {
					// Create full destination
					address := fmt.Sprintf("%s:%d", ip, port)

					// Set a shgort timeout
					timeout := 2 * time.Second
					conn, err := net.DialTimeout("tcp", address, timeout)
					if err == nil {
						hostActive = true
						activePorts = append(activePorts, port)
						conn.Close()
					}

				}

				// After host scan has completed, if host had any ports open
				if hostActive {
					scanResults.ActiveHosts[ip] = activePorts
					log.Info().Str("ip", ip).
						Msg(fmt.Sprintf("Open Ports: %v", activePorts))
					Database.InsertPortScan(ip, activePorts, time.Now())
				} else {
					if !slices.Contains(scanResults.InactiveHosts, ip) {
						scanResults.InactiveHosts = append(
							scanResults.InactiveHosts,
							ip,
						)
						// log.Debug().Str("ip", ip).Msg("No ports found.")
					}
				}
			}(addr.String())
		}
	}

	wg.Wait()
}
