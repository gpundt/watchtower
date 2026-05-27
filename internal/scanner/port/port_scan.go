package port

import (
	"net"
	"slices"
	"sync"
	"time"

	Config "watchtower/internal/config"

	"github.com/rs/zerolog/log"
)

type PortScanResults struct {
	ActiveHosts map[]string[]int
	InactiveHosts []string
}

func RunPortScan(subnets []string) {
	log.Info().Msg("Beginning ARP Scan")

	scanResults := PortScanResults{
		ActiveHosts := map[]string{},
		InactiveHosts := []string{},
	}
	guard := make(
		chan struct{},
		Config.ServerConfig.Scanner.MaxConcurrentScans,
	)

	var wg sync.WaitGroup
	
	// Iterate over each subnet in subnets
	for _, subnet := range subnets {
		_, ipv4Net, _ := net.ParseCIDR(subnet)
		mask := binary.BigEndian.Uint32(ipv4Net.Mask)
		networkStart := binary.BigEndian.Uint32(ipv4Net.IP)
		networkEnd := (networkStart & mask) | ^mask

		for i := networkStart; i <= networkEnd; i++ {
			emptyIP := make(net.IP, 4)
			binary.BigEndian.PutUint32(emptyIP, i)
			hostIP := emptyIP.String()
			
			wg.Add(1)
			guard <- struct{}{} // Block if max number of goroutines alrady running

			go func(ip string) {
				defer wg.Done()
				defer func() { <-guard }() // Release the slot in the semaphore

				hostActive := false
				activePorts := []int{}
				for _, port := range Config.ServerConfig.Scanner.Ports {
					// Create full destination
					address := net.JoinHostPort(ip, port)

					// Set a shgort timeout
					timeout := 2 * time.Second
					conn, err := net.DialTimeout("tcp", address, timeout)
					if err == nil {
						hostActive = true
						activePorts = append(activePorts, port)
					}
					conn.Close()
				}

				// After host scan has completed, if host had any ports open
				if hostActive {
					scanResults.ActiveHosts[ip] = activePorts
				} else {
					if !slices.Contains(scanResults.InactiveHosts, ip) {
						scanResults.InactiveHosts = append(
							scanResults.InactiveHosts,
							ip,
						)
					}
				}
			}(hostIP)
		}
	}

	wg.Wait()
}