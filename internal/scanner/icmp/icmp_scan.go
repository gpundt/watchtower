package icmp


import (
	"encoding/binary"
	"sync"
	"time"
	"net"

	Config "watchtower/internal/config"

	"github.com/prometheus-community/pro-bing"
	"github.com/rs/zerolog/log"
)

type ICMPScanResults struct {
	ActiveHosts		[]string
	InactiveHosts	[]string
}

func RunICMPScan(subnets []string) {
	log.Info().Msg("Beginning ICMP Scan")
	
	scanResults := ICMPScanResults{
		ActiveHosts: []string{},
		InactiveHosts: []string{},
	}
	guard := make(chan struct{}, Config.ServerConfig.Scanner.MaxConcurrentScans)

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
				defer func() { <-guard }() // Release the slot in the semphore
			
				pinger, err := probing.NewPinger(hostIP)
				if err != nil {
					log.Err(err).Str("hostIP", hostIP).Msg("")
					scanResults.InactiveHosts = append(scanResults.InactiveHosts, hostIP)
					return
				}	

				pinger.Count = 2
				pinger.Timeout = 2 * time.Second

				// Run blocks until the probes are complete or timeout is reached
				err = pinger.Run()
				stats := pinger.Statistics()

				if err != nil || stats.PacketsRecv == 0 {
					log.Debug().Str("hostIP", hostIP).Msg("Probe Missed")
					scanResults.InactiveHosts = append(scanResults.InactiveHosts, hostIP)
				} else {
					log.Debug().Str("hostIP", hostIP).Msg("Probe Hit!")
					scanResults.ActiveHosts = append(scanResults.ActiveHosts, hostIP)
				}
			}(hostIP)
		}
	}

	wg.Wait()
}