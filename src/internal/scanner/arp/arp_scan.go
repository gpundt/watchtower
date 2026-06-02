package arp

import (
	"bytes"
	"context"
	"encoding/binary"
	"net"
	"time"

	Database "watchtower/internal/database"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
	"github.com/rs/zerolog/log"
)

type ARPScanResults struct {
	ActiveHosts   []string
	InactiveHosts []string
}

func RunARPScan(subnets map[string]string) {
	log.Info().Msg("Beginning ARP Scan")

	// For each interface and subnet
	for interfaceName := range subnets {

		iface, err := net.InterfaceByName(interfaceName)
		if err != nil {
			log.Err(err).Str("func", "RunARPScan").Msg("")
			return
		}

		var ipNet *net.IPNet
		addrs, _ := iface.Addrs()
		for _, addr := range addrs {
			if ip, ok := addr.(*net.IPNet); ok && ip.IP.To4() != nil {
				ipNet = ip
				break
			}
		}

		// Open network interface for live packet capture
		handle, err := pcap.OpenLive(interfaceName, 65536, true, pcap.BlockForever)
		if err != nil {
			log.Err(err).Str("func", "RunARPScan").Msg("")
			return
		}
		defer handle.Close()

		// Context to control execution timeout
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		// Start Listeneing for incoming arp replies in a goroutine
		go listenForARPReplies(ctx, handle, iface)

		// Broadcast ARP request across local subnet
		if err := generateARPRequest(handle, iface, ipNet); err != nil {
			log.Err(err).Str("func", "RunARPScan").Msg("Failed to send ARP requests")
			return
		}

		//Allow remaining responses to arrive before exiting
		<-ctx.Done()
	}
}

// Function to filter and log incoming frames containing ARP answers
func listenForARPReplies(
	ctx context.Context,
	handle *pcap.Handle,
	iface *net.Interface,
) {
	src := gopacket.NewPacketSource(handle, layers.LayerTypeEthernet)
	in := src.Packets()

	for {
		select {
		case <-ctx.Done():
			return
		case packet := <-in:
			arpLayer := packet.Layer(layers.LayerTypeARP)
			if arpLayer == nil {
				continue
			}
			arp := arpLayer.(*layers.ARP)

			// Only process replies targeting this server
			if arp.Operation == layers.ARPReply && bytes.Equal(arp.DstHwAddress, iface.HardwareAddr) {
				ip := net.IP(arp.SourceProtAddress)
				log.Debug().Str("ip", ip.String()).
					Msg("ARP Reply Received!")
				Database.InsertPortScan(ip.String(), nil, time.Now())

			}
		}
	}
}

// Loops through all IP addresses in the subnet range and sends ARP broadcast
func generateARPRequest(
	handle *pcap.Handle,
	iface *net.Interface,
	ipNet *net.IPNet,
) error {
	// Base Ethernet layer
	eth := layers.Ethernet{
		SrcMAC:       iface.HardwareAddr,
		DstMAC:       net.HardwareAddr{0xff, 0xff, 0xff, 0xff, 0xff, 0xff},
		EthernetType: layers.EthernetTypeARP,
	}

	//Base ARP request layer template
	arp := layers.ARP{
		AddrType:          layers.LinkTypeEthernet,
		Protocol:          layers.EthernetTypeIPv4,
		HwAddressSize:     6,
		ProtAddressSize:   4,
		Operation:         layers.ARPRequest,
		SourceHwAddress:   []byte(iface.HardwareAddr),
		SourceProtAddress: []byte(ipNet.IP.To4()),
		DstHwAddress:      []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
	}

	buf := gopacket.NewSerializeBuffer()
	opts := gopacket.SerializeOptions{FixLengths: true, ComputeChecksums: true}

	// Calculate start IP and end IP for subnet
	mask := binary.BigEndian.Uint32(ipNet.Mask)
	start := binary.BigEndian.Uint32(ipNet.IP.To4()) & mask
	end := start | ^mask

	// Iterate and send requests to every valid host within subnet
	for i := start + 1; i < end; i++ {
		targetIP := make(net.IP, 4)
		binary.BigEndian.PutUint32(targetIP, i)
		arp.DstProtAddress = []byte(targetIP.To4())

		if err := gopacket.SerializeLayers(buf, opts, &eth, &arp); err != nil {
			return err
		}
		if err := handle.WritePacketData(buf.Bytes()); err != nil {
			return err
		}
		buf.Clear()
		time.Sleep(2 * time.Millisecond)
	}

	return nil
}
