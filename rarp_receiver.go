package main

import (
	"encoding/hex"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"syscall"
	"time"
)

const (
	RARPEtherType = 0x8035
)

func main() {
	if len(os.Args) != 2 {
		log.Fatalf("Usage: %s <network-interface>", os.Args[0])
	}

	interfaceName := os.Args[1]

	// Create a raw socket
	fd, err := syscall.Socket(syscall.AF_PACKET, syscall.SOCK_RAW, int(htons(RARPEtherType)))
	if err != nil {
		log.Fatalf("Failed to create socket: %v", err)
	}
	defer syscall.Close(fd)

	// Bind the socket to the network interface
	iface, err := net.InterfaceByName(interfaceName)
	if err != nil {
		log.Fatalf("Failed to get interface: %v", err)
	}

	sll := syscall.SockaddrLinklayer{
		Protocol: htons(RARPEtherType),
		Ifindex:  iface.Index,
	}
	if err := syscall.Bind(fd, &sll); err != nil {
		log.Fatalf("Failed to bind socket: %v", err)
	}

	// Receive packets
	buf := make([]byte, 65536)
	for {
		numbytes, _, err := syscall.Recvfrom(fd, buf, 0)
		if err != nil {
			log.Fatalf("Failed to receive packet: %v", err)
		}

		printPacket(buf[:numbytes])
	}
}

func htons(i uint16) uint16 {
	return (i<<8)&0xff00 | i>>8
}

func printPacket(packet []byte) {
	timestamp := time.Now().Format("15:04:05.000000")
	fmt.Printf("%s %s\n", timestamp, formatPacketData(packet))
}

func formatPacketData(data []byte) string {
	const bytesPerLine = 16
	lines := make([]string, 0, (len(data)+bytesPerLine-1)/bytesPerLine)

	for i := 0; i < len(data); i += bytesPerLine {
		end := i + bytesPerLine
		if end > len(data) {
			end = len(data)
		}
		line := fmt.Sprintf("%04x  %-47s  %s",
			i,
			hex.EncodeToString(data[i:end]),
			formatASCII(data[i:end]),
		)
		lines = append(lines, line)
	}

	return strings.Join(lines, "\n")
}

func formatASCII(data []byte) string {
	const printableMin = 32
	const printableMax = 126

	var result strings.Builder
	for _, b := range data {
		if b >= printableMin && b <= printableMax {
			result.WriteByte(b)
		} else {
			result.WriteByte('.')
		}
	}

	return result.String()
}
