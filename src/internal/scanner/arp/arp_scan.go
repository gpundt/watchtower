package arp

import (
	"github.com/rs/zerolog/log"
)

type ARPScanResults struct {
	ActiveHosts   []string
	InactiveHosts []string
}

func RunARPScan() {
	log.Info().Msg("Beginning ARP Scan")
}
