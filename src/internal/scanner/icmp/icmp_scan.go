package icmp

import (
	"net/netip"
	"sync"
	"time"

	Config "watchtower/internal/config"
	Database "watchtower/internal/database"

	Probing "github.com/prometheus-community/pro-bing"
	"github.com/rs/zerolog/log"
)

type ICMPScanResults struct {
	ActiveHosts   []string
	InactiveHosts []string
}

// Function to run ICMP scan on every /24 subnet he server is on
func RunICMPScan(subnets []string) {
	log.Info().Msg("Beginning ICMP Scan")

	scanResults := ICMPScanResults{
		ActiveHosts:   []string{},
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
				defer func() { <-guard }() // Release the slot in the semphore

				// Create object to ping hosts
				pinger, err := Probing.NewPinger(ip)
				if err != nil {
					log.Err(err).Str("ip", ip).Msg("")
					scanResults.InactiveHosts = append(
						scanResults.InactiveHosts,
						ip,
					)
					return
				}

				pinger.Count = 2
				pinger.Timeout = 2 * time.Second

				// Run blocks until the probes are complete or timeout is reached
				err = pinger.Run()
				stats := pinger.Statistics()

				// If error or no response, we missed
				if err != nil || stats.PacketsRecv == 0 {
					// log.Debug().Str("ip", ip).Msg("ICMP Probe Missed")
					scanResults.InactiveHosts = append(
						scanResults.InactiveHosts,
						ip,
					)
				} else {
					log.Debug().Str("ip", ip).Msg("ICMP Probe Hit!")
					scanResults.ActiveHosts = append(
						scanResults.ActiveHosts,
						ip,
					)
					Database.InsertPortScan(ip, nil, time.Now())
				}
			}(addr.String())
		}
	}

	wg.Wait()
}
